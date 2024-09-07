package services

import (
	"github.com/cantylv/thumbnail-loader/services/connectors"
	"go.uber.org/zap"
)

type Services struct {
	DBCacheClient       connectors.ClientDB
	InMemoryCacheClient connectors.ClientCache
	MinioClient         connectors.ClientS3
}

// var _ sqlite.Client = sqlite.Init(zap.Must(zap.NewProduction()))

func Init(logger *zap.Logger, inMemoryClient connectors.EngineCache, dbCacheClient connectors.EngineDB, minioClient connectors.EngineS3) (cluster *Services) {
	return &Services{
		InMemoryCacheClient: inMemoryClient.InitClientCache(logger),
		DBCacheClient:       dbCacheClient.InitClientDB(logger),
		MinioClient:         minioClient.InitClientS3(logger),
	}
}
