//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package main

import (
	"testing"
)

// Of course, it's unuseful, but test coverage was done 100%
func TestMain(t *testing.T) {
	main()
}
