package services

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/thumbnail-loader/services/memcached"
	minios3 "github.com/cantylv/thumbnail-loader/services/minio"
	minio "github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Services struct {
	CacheClient *memcache.Client
	MinioClient *minio.Client
}

func Init(logger *zap.Logger) *Services {
	return &Services{
		CacheClient: memcached.Init(logger),
		MinioClient: minios3.Init(logger),
	}
}
