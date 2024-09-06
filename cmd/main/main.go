package main

import (
	"github.com/cantylv/thumbnail-loader/config"
	"github.com/cantylv/thumbnail-loader/internal/app"
	"go.uber.org/zap"
)

func main(){
	logger := zap.Must(zap.NewProduction())
	// reading env variables, config file, cli flags
	config.Read("./config/config.yaml", logger)
	// initalize app instance
	app.Run(logger)
	logger.Info("app has successfully completed its work")
}
