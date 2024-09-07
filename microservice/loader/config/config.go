package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func setDefaultParameters() {
	// __memcached__ variables
	viper.SetDefault("memcached.host", "localhost")
	viper.SetDefault("memcached.port", 11211)
	viper.SetDefault("memcached.cache_timeout", time.Minute*10)

	// __grpc_download variables
	viper.SetDefault("grpc_loader.host", "localhost")
	viper.SetDefault("grpc_loader.port", 8000)

	// __minio__ variables
	viper.SetDefault("minio.host", "localhost")
	viper.SetDefault("minio.port", 9000)
	viper.SetDefault("minio.bucket_name", "images")
	viper.SetDefault("minio.access_key", "admin")
	viper.SetDefault("minio.secret_access_key", 123)
	viper.SetDefault("minio.use_ssl", false)
}

func Read(path string, logger *zap.Logger) {
	setDefaultParameters()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Panic(fmt.Sprintf("fatal error config file: %v", err))
		}
		logger.Warn("warning: configuration file is not found, programm will be executed within default configuration")
	}
	logger.Info("successful read of configuration")
}
