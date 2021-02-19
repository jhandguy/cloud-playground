compile:
	protoc --go_out=. --go-grpc_out=. --proto_path=../dynamo proto/item/item.pg.proto
	protoc --go_out=. --go-grpc_out=. --proto_path=../s3 proto/object/object.pg.proto

lint:
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

PORT ?= 8080
test: export GATEWAY_PORT=$(PORT)
test: export GATEWAY_API_KEY=$(shell terraform -chdir=../terraform output -json gateway_api_key | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=../terraform output -json gateway_url | jq -r .)
test: export S3_URL=$(shell terraform -chdir=../terraform output -json s3_url | jq -r .)
test: export S3_TOKEN=$(shell terraform -chdir=../terraform output -json s3_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=../terraform output -json dynamo_url | jq -r .)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=../terraform output -json dynamo_token | jq -r .)
test:
	go test ./... -cover -race