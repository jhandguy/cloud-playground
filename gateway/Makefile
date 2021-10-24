compile:
	protoc --go_out=. --go-grpc_out=. --proto_path=../dynamo proto/item/item.pg.proto
	protoc --go_out=. --go-grpc_out=. --proto_path=../s3 proto/object/object.pg.proto

lint:
	helm lint helm
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

HTTP_PORT ?= 8080
METRICS_PORT ?= 9090
test: export GATEWAY_HTTP_PORT=$(HTTP_PORT)
test: export GATEWAY_METRICS_PORT=$(METRICS_PORT)
test: export GATEWAY_TOKEN=$(shell terraform -chdir=../terraform output -json gateway_token | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=../terraform output -json ingress_gateway_url | jq -r .)
test: export S3_URL=$(shell terraform -chdir=../terraform output -json s3_url | jq -r .)
test: export S3_TOKEN=$(shell terraform -chdir=../terraform output -json s3_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=../terraform output -json dynamo_url | jq -r .)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=../terraform output -json dynamo_token | jq -r .)
test:
	go test ./... -cover -race

load: export GATEWAY_TOKEN=$(shell terraform -chdir=../terraform output -json gateway_token | jq -r .)
load: export GATEWAY_URL=$(shell terraform -chdir=../terraform output -json ingress_gateway_url | jq -r .)
load:
	k6 run k6/script.js

update:
	go get -u
	go mod tidy

docker:
	docker buildx build --tag ghcr.io/jhandguy/devops-playground/gateway:latest --tag ghcr.io/jhandguy/devops-playground/gateway:1.0.0 --file Dockerfile . --push
