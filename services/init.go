package services

import (
	"database/sql"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/thumbnail-loader/services/memcached"
	minios3 "github.com/cantylv/thumbnail-loader/services/minio"
	"github.com/cantylv/thumbnail-loader/services/sqlite"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Services struct {
	DBCacheClient       *sql.DB
	InMemoryCacheClient *memcache.Client
	MinioClient         *minio.Client
}

func Init(logger *zap.Logger) (cluster *Services) {
	cluster = new(Services)
	if viper.GetBool("cache_inmemory") {
		cluster.InMemoryCacheClient = memcached.Init(logger)
	} else {
		cluster.DBCacheClient = sqlite.Init(logger)
	}
	cluster.MinioClient = minios3.Init(logger)
	return cluster
}
