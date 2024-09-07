package memcached

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/thumbnail-loader/services/connectors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ClientInstance struct{}

var _ connectors.EngineCache = (*ClientInstance)(nil)

func NewClientInstanse() *ClientInstance {
	return &ClientInstance{}
}

func (t *ClientInstance) InitClientCache(logger *zap.Logger) connectors.ClientCache {
	connLine := fmt.Sprintf("%s:%d", viper.GetString("memcached.host"), viper.GetUint16("memcached.port"))
	client := memcache.New(connLine)
	client.MaxIdleConns = 20

	logger.Info("succesful connection to Memcached")
	return client
}
