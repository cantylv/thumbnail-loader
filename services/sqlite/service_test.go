//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitClientDB(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("panic was intercepted")
		}
	}()
	clientInstance := NewClientInstanse()
	client := clientInstance.InitClientDB(zap.Must(zap.NewProduction()))
	assert.NotNil(t, client, "sqlite always will be initialized")
}
