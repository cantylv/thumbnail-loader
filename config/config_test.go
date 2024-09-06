//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError) // сбрасываем флаги
}

func setTestEnv() {
	viper.Reset()
	os.Args = []string{"main"}
	os.Unsetenv("UPLOAD_FOLDER")
}

// In this test we specially set wrong path to the configuration file
func TestDefaultConfiguration(t *testing.T) {
	setTestEnv()
	Read("./blabla.yaml", zap.Must(zap.NewProduction()))

	assert.Equal(t, false, viper.GetBool("cache_inmemory"), "default value for viper variable 'cache_inmemory' is false")
	assert.Equal(t, false, viper.GetBool("async"), "default value for viper variable 'async' is false")
	assert.Equal(t, time.Second*30, viper.GetDuration("cache_timeout"), "default value for viper variable 'cache_timeout' is 30s")
	assert.Equal(t, "uploads", viper.GetString("upload_folder"), "default value for viper variable 'upload_folder' is 'uploads'")
	assert.Equal(t, "localhost", viper.GetString("grpc_loader.host"), "default value for viper variable 'grpc_loader.host' is 'localhost'")
	assert.Equal(t, 8000, viper.GetInt("grpc_loader.port"), "default value for viper variable 'grpc_loader.host' is 8000")

	resetFlags()
}

func TestCustomConfigurationUsingYAML(t *testing.T) {
	setTestEnv()
	Read("./config.yaml", zap.Must(zap.NewProduction()))

	assert.Equal(t, true, viper.GetBool("cache_inmemory"), "value for viper variable 'cache_inmemory' from config.yaml is true")
	assert.Equal(t, true, viper.GetBool("async"), "value for viper variable 'async' from config.yaml is true")
	assert.Equal(t, time.Second*60, viper.GetDuration("cache_timeout"), "value for viper variable 'cache_timeout' from config.yaml is 60s")
	assert.Equal(t, "thumbnails", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from config.yaml is 'thumbnails'")
	assert.Equal(t, "localhost", viper.GetString("grpc_loader.host"), "value for viper variable 'grpc_loader.host' from config.yaml is 'localhost'")
	assert.Equal(t, 8001, viper.GetInt("grpc_loader.port"), "value for viper variable 'grpc_loader.host' from config.yaml is 8001")

	resetFlags()
}

func TestEnvVar(t *testing.T) {
	setTestEnv()
	Read("./config.yaml", zap.Must(zap.NewProduction()))

	assert.Equal(t, "thumbnails", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from config.yaml is 'thumbnails'")
	err := os.Setenv("UPLOAD_FOLDER", "youtube")
	if err != nil {
		t.Fatalf("error while setting env variable: %v", err)
	}
	assert.Equal(t, "youtube", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from env is 'youtube', it has more priority than conf file")

	resetFlags()
}

func TestEnvVarBeforeReadConfig(t *testing.T) {
	setTestEnv()
	err := os.Setenv("UPLOAD_FOLDER", "youtube")
	if err != nil {
		t.Fatalf("error while setting env variable: %v", err)
	}
	Read("./config.yaml", zap.Must(zap.NewProduction()))
	assert.Equal(t, "youtube", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from env is 'youtube', it has more priority than conf file")

	resetFlags()
}

func TestCmdFlags(t *testing.T) {
	setTestEnv()
	Read("./config.yaml", zap.Must(zap.NewProduction()))

	assert.Equal(t, "thumbnails", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from config.yaml is 'thumbnails'")
	err := os.Setenv("UPLOAD_FOLDER", "youtube")
	if err != nil {
		t.Fatalf("error while setting env variable: %v", err)
	}
	assert.Equal(t, "youtube", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from env is 'youtube', it has more priority than conf file")

	resetFlags()

	var uploadFolder string
	pflag.StringVarP(&uploadFolder, "upload_folder", "u", "uploads",
		"the destination folder for uploading files from youtube")
	os.Args = []string{"cmd", "--upload_folder=images"}
	pflag.Parse()

	viper.BindPFlag("upload_folder", pflag.Lookup("upload_folder"))
	assert.Equal(t, "images", viper.GetString("upload_folder"), "value for viper variable 'upload_folder' from cmd line is 'images', it has more priority than env var")
	resetFlags()
}
