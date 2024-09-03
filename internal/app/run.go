package app

import (
	"fmt"

	"github.com/cantylv/thumbnail-loader/internal/utils/functions"
	"github.com/cantylv/thumbnail-loader/services"
	"go.uber.org/zap"
)

// Run
// start app engine (logic)
func Run(logger *zap.Logger) {
	// initialization of rdbms, s3, in-memory storage
	serviceCluster := services.Init(logger)
	defer func(cluster *services.Services) {
		err := serviceCluster.CacheClient.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error while closing memcached: %v", err))
		}
	}(serviceCluster)

	functions.StartEngine(serviceCluster, logger)
}
