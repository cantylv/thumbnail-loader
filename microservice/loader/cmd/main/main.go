package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/cantylv/thumbnail-loader/config"
	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/repository/urls"
	ucUrls "github.com/cantylv/thumbnail-loader/microservice/loader/internal/usecase/urls"
	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	"github.com/cantylv/thumbnail-loader/services"
	"github.com/cantylv/thumbnail-loader/services/memcached"
	"github.com/cantylv/thumbnail-loader/services/minio"
	"github.com/cantylv/thumbnail-loader/services/sqlite"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	config.Read("./microservice/loader/config/config.yaml", logger)

	address := fmt.Sprintf("%s:%d", viper.GetString("grpc_loader.host"), viper.GetInt("grpc_loader.port"))
	conn, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("microservice \"download\" doesn't respond", zap.String("error", err.Error()))
	}
	logger.Info(fmt.Sprintf("microservice \"download\" responds on address %s", address))

	// init grpc server
	server := grpc.NewServer()
	// initialization of rdbms, s3, in-memory storage
	inMemoryClient := memcached.NewClientInstanse()
	s3Client := minio.NewClientInstanse()
	dbClient := sqlite.NewClientInstanse()
	serviceCluster := services.Init(logger, inMemoryClient, dbClient, s3Client)
	defer func(cluster *services.Services) {
		if serviceCluster.InMemoryCacheClient != nil {
			err := serviceCluster.InMemoryCacheClient.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("error while closing memcached: %v", err))
			}
		}
		if serviceCluster.DBCacheClient != nil {
			err := serviceCluster.DBCacheClient.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("error while closing sqlite: %v", err))
			}
		}
	}(serviceCluster)
	repoLayer := urls.NewRepoLayer(serviceCluster.DBCacheClient)
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = repoLayer.Init(context)
	if err != nil {
		logger.Fatal(err.Error())
	}
	usecaseLayer := ucUrls.NewUsecaseLayer(repoLayer, serviceCluster, logger)

	gen.RegisterDownloadManagerServer(server, usecaseLayer)
	err = server.Serve(conn)
	if err != nil {
		logger.Fatal(err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.GracefulStop()
	logger.Info("microservice \"download\" has shut down")
	os.Exit(0)
}
