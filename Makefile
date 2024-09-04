lint:
	goimports -e ./ 
	go vet ./...
	golangci-lint run ./...

easyjs:
	easyjson -no_std_marshalers -all internal/entity

pb:
	protoc -I microservice/loader/proto microservice/loader/proto/loader.proto --go_out=microservice/loader/proto --go-grpc_out=microservice/loader/proto

vendor:
	go mod vendor

run: 
	docker-compose up -d