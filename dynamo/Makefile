compile:
	protoc --go_out=. --go-grpc_out=. proto/item/item.pg.proto

lint:
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

PORT ?= 8080
test: export AWS_DYNAMO_ENDPOINT=$(shell terraform -chdir=../terraform output -json aws_dynamo_endpoint | jq -r .)
test: export AWS_DYNAMO_TABLE=$(shell terraform -chdir=../terraform output -json aws_dynamo_table | jq -r .)
test: export DYNAMO_PORT=$(PORT)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=../terraform output -json dynamo_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=../terraform output -json dynamo_url | jq -r .)
test:
	go test ./... -cover -race