lint:
	goimports -e ./ 
	go vet ./...
	golangci-lint run ./...

pb:
	protoc -I microservice/loader/proto microservice/loader/proto/loader.proto --go_out=microservice/loader/proto --go-grpc_out=microservice/loader/proto

vendor:
	go mod vendor

tidy:
	go mod tidy

mock:
	mockgen -source=microservice/loader/proto/gen/loader_grpc.pb.go -destination=microservice/loader/mocks/mock_download_manager_client.go -package=mocks 
	mockgen -source=services/connectors/define.go -destination=services/mocks/mock_clients.go -package=mocks 
	
gen:
	go generate ./...

init: gen tidy pb 
	mkdir ./services/minio/data ./services/sqlite/data
	touch ./services/sqlite/data/database.db

run: vendor
	docker-compose up -d
	go run ./microservice/loader/cmd/main