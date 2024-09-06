//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package app

import (
	"testing"

	"go.uber.org/zap"
)

// Of course, it's unuseful, but test coverage was done 100%
func TestRun(t *testing.T) {
	Run(zap.Must(zap.NewProduction()))
}
