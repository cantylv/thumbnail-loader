//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package main

import "testing"

func TestMain(t *testing.T) {
	main()
}
