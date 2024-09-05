package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type CmdFlags struct {
	NeedCache    bool
	Async        bool
	CacheTimeout time.Duration
	UploadFolder string
}

func setDefaultParameters() {
	viper.AutomaticEnv()

	// __memcached__ variables
	viper.SetDefault("memcached.host", "memcached")
	viper.SetDefault("memcached.port", 11211)
	viper.SetDefault("memcached.cache_timeout", time.Minute*10)

	// __grpc_download variables
	viper.SetDefault("grpc_loader.host", "localhost")
	viper.SetDefault("grpc_loader.port", 8001)

	// __project__ variables
	viper.SetDefault("async", false)
	viper.SetDefault("cache_inmemory", false)

	// __minio__ variables
	viper.SetDefault("minio.host", "minio")
	viper.SetDefault("minio.port", 9000)
	viper.SetDefault("minio.bucket_name", "images")
	viper.SetDefault("minio.access_key", "admin")
	viper.SetDefault("minio.secret_access_key", 123)
	viper.SetDefault("minio.use_ssl", false)

	if viper.Get("UPLOAD_FOLDER") == nil {
		viper.SetDefault("upload_folder", "uploads")
	} else {
		viper.Set("upload_folder", viper.Get("UPLOAD_FOLDER"))
	}
}

func readCmdFlags() {
	var needCache bool
	var async bool
	var cacheTimeout time.Duration
	var uploadFolder string

	pflag.BoolVarP(&needCache, "cache_inmemory", "c", false,
		"determines 'type' of cache; if true, cache data will be stored in ram, in another way in winchester")
	pflag.BoolVarP(&async, "async", "a", false,
		"configure whether asynchronous loading is required")
	pflag.DurationVarP(&cacheTimeout, "cache_timeout", "t", time.Minute*10,
		"the duration for which cache instance will store data")
	pflag.StringVarP(&uploadFolder, "upload_folder", "u", "uploads",
		"the destination folder for uploading files from youtube")
	pflag.Parse()

	viper.BindPFlag("cache_inmemory", pflag.Lookup("cache_inmemory"))
	viper.BindPFlag("async", pflag.Lookup("async"))
	viper.BindPFlag("memcached.cache_timeout", pflag.Lookup("cache_timeout"))
	viper.BindPFlag("upload_folder", pflag.Lookup("upload_folder"))
}

func Read(path string, logger *zap.Logger) {
	setDefaultParameters()
	readCmdFlags()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			logger.Panic(fmt.Sprintf("fatal error config file: %v", err))
		}
		logger.Warn("warning: configuration file is not found, programm will be executed within default configuration")
	}
	logger.Info("successful read of configuration")
}
