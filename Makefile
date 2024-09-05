lint:
	goimports -e ./ 
	go vet ./...
	golangci-lint run ./...

easyjs:
	easyjson -no_std_marshalers -all microservice/loader/internal/entity

pb:
	protoc -I microservice/loader/proto microservice/loader/proto/loader.proto --go_out=microservice/loader/proto --go-grpc_out=microservice/loader/proto

vendor:
	go mod vendor

tidy:
	go mod tidy

init: tidy easyjs pb 
	mkdir ./services/minio/data ./services/sqlite/data
	touch ./services/sqlite/data/database.db

run: vendor
	docker-compose up -d
	go run ./microservice/loader/cmd/main