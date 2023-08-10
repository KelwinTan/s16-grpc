dep: ## for updating go mod imports and vendor
	go mod tidy && go mod vendor

run: ## for running binary
	go run cmd/omdb/main.go

gen-proto: ## for generating protobufs
	protoc --go_out=./omdb --go_opt=paths=source_relative \
    --go-grpc_out=./omdb --go-grpc_opt=paths=source_relative \
    omdb.proto

mock:
	mockgen -source=./api/proto/v1/omdb/omdb.pb.go OmdbServiceClient