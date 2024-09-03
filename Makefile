lint:
	goimports -e ./ 
	go vet ./...
	golangci-lint run ./...

easyjs:
	easyjson -no_std_marshalers -all internal/entity