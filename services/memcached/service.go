package memcached

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *memcache.Client {
	connLine := fmt.Sprintf("%s:%d", viper.GetString("memcached.host"), viper.GetUint16("memcached.port"))
	client := memcache.New(connLine)
	client.MaxIdleConns = 20

	logger.Info("succesful connection to Memcached")
	return client
}
