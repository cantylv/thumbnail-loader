lint:
	goimports -e ./ 
	go vet ./...
	golangci-lint run ./...

easyjs:
	easyjson -no_std_marshalers -all microservice/loader/internal/entity

pb:
	protoc -I microservice/loader/proto microservice/loader/proto/loader.proto --go_out=microservice/loader/proto --go-grpc_out=microservice/loader/proto

vendor:
	go mod tidy
	go mod vendor

run: vendor
	docker-compose up -d
	go run ./microservice/loader/cmd/main