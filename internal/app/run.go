package app

import (
	"fmt"

	"github.com/cantylv/thumbnail-loader/internal/functions"
	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run
// start app engine (logic)
func Run(logger *zap.Logger) {
	// init grpc client
	serverConnect, err := grpc.NewClient(fmt.Sprintf("%s:%d", viper.GetString("grpc_loader.host"), viper.GetInt("grpc_loader.port")),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func(serverConnect *grpc.ClientConn) {
		err := serverConnect.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(serverConnect)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error while creating grpc client: %v", err))
	}
	err = functions.StartEngine(gen.NewDownloadManagerClient(serverConnect), logger)
	if err != nil {
		logger.Error(err.Error())
	}
}
