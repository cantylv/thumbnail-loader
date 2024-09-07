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

	// __memcached__ variables
	assert.Equal(t, "localhost", viper.GetString("memcached.host"), "default value for viper variable 'memcached.host' is 'localhost'")
	assert.Equal(t, 11211, viper.GetInt("memcached.port"), "default value for viper variable 'memcached.port' is 8000")
	assert.Equal(t, 10*time.Minute, viper.GetDuration("memcached.cache_timeout"), "default value for viper variable 'memcached.cache_timeout' is 10m")
	// __grpc_download variables
	assert.Equal(t, "localhost", viper.GetString("grpc_loader.host"), "default value for viper variable 'grpc_loader.host' is 'localhost'")
	assert.Equal(t, 8000, viper.GetInt("grpc_loader.port"), "default value for viper variable 'grpc_loader.port' is 8000")
	// __minio__ variables
	assert.Equal(t, "localhost", viper.GetString("minio.host"), "default value for viper variable 'minio.host' is 'localhost'")
	assert.Equal(t, 9000, viper.GetInt("minio.port"), "default value for viper variable 'minio.port' is 9000")
	assert.Equal(t, "images", viper.GetString("minio.bucket_name"), "default value for viper variable 'minio.bucket_name' is 'images'")
	assert.Equal(t, "admin", viper.GetString("minio.access_key"), "default value for viper variable 'minio.access_key' is 'admin'")
	assert.Equal(t, 123, viper.GetInt("minio.secret_access_key"), "default value for viper variable 'minio.secret_access_key' is 123")
	assert.Equal(t, false, viper.GetBool("minio.use_ssl"), "default value for viper variable 'minio.use_ssl' is false")
	resetFlags()
}

func TestCustomConfigurationUsingYAML(t *testing.T) {
	setTestEnv()
	Read("./config.yaml", zap.Must(zap.NewProduction()))

	// __memcached__ variables
	assert.Equal(t, "localhost", viper.GetString("memcached.host"), "default value for viper variable 'memcached.host' is 'localhost'")
	assert.Equal(t, 22122, viper.GetInt("memcached.port"), "default value for viper variable 'memcached.port' is 8000")
	assert.Equal(t, 30*time.Minute, viper.GetDuration("memcached.cache_timeout"), "default value for viper variable 'memcached.cache_timeout' is 10m")
	// __grpc_download variables
	assert.Equal(t, "localhost", viper.GetString("grpc_loader.host"), "default value for viper variable 'grpc_loader.host' is 'localhost'")
	assert.Equal(t, 8001, viper.GetInt("grpc_loader.port"), "default value for viper variable 'grpc_loader.port' is 8000")
	// __minio__ variables
	assert.Equal(t, "localhost", viper.GetString("minio.host"), "default value for viper variable 'minio.host' is 'localhost'")
	assert.Equal(t, 9000, viper.GetInt("minio.port"), "default value for viper variable 'minio.port' is 9000")
	assert.Equal(t, "thumbnails", viper.GetString("minio.bucket_name"), "default value for viper variable 'minio.bucket_name' is 'images'")
	assert.Equal(t, "admin", viper.GetString("minio.access_key"), "default value for viper variable 'minio.access_key' is 'admin'")
	assert.Equal(t, "admin123", viper.GetString("minio.secret_access_key"), "default value for viper variable 'minio.secret_access_key' is 'admin123'")
	assert.Equal(t, false, viper.GetBool("minio.use_ssl"), "default value for viper variable 'minio.use_ssl' is false")

	resetFlags()
}
