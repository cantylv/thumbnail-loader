//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package memcached

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitClientCache(t *testing.T) {
	clientInstance := NewClientInstanse()
	client := clientInstance.InitClientCache(zap.Must(zap.NewProduction()))
	assert.NotNil(t, client, "memcached always returns not nil object")
}
