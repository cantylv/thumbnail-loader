package services

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	repoUrls "github.com/cantylv/thumbnail-loader/internal/repository/urls"
	ucUrls "github.com/cantylv/thumbnail-loader/internal/usecase/urls"
	"github.com/cantylv/thumbnail-loader/services/memcached"
	minios3 "github.com/cantylv/thumbnail-loader/services/minio"
	"github.com/cantylv/thumbnail-loader/services/sqlite"
	minio "github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Services struct {
	DBCacheClient       ucUrls.Usecase
	InMemoryCacheClient *memcache.Client
	MinioClient         *minio.Client
}

func Init(logger *zap.Logger) (cluster *Services) {
	cluster = new(Services)
	if viper.GetBool("cache_inmemory") {
		cluster.InMemoryCacheClient = memcached.Init(logger)
	} else {
		repoLayer := repoUrls.NewRepoLayer(sqlite.Init(logger))
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := repoLayer.Init(context)
		if err != nil {
			logger.Error(err.Error())
		}
		cluster.DBCacheClient = ucUrls.NewUsecaseLayer(repoLayer)
	}
	cluster.MinioClient = minios3.Init(logger)
	return cluster
}
