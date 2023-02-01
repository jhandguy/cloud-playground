ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
HTTP_PORT ?= 8080
METRICS_PORT ?= 9090

compile:
	protoc --go_out=. --go-grpc_out=. --proto_path=../dynamo proto/item/item.pg.proto
	protoc --go_out=. --go-grpc_out=. --proto_path=../s3 proto/object/object.pg.proto

lint_helm:
	helm lint helm

lint:
	golangci-lint run

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
	go test ./... -short -cover -race

test: export GATEWAY_HTTP_PORT=$(HTTP_PORT)
test: export GATEWAY_METRICS_PORT=$(METRICS_PORT)
test: export GATEWAY_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json gateway_token | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_url | jq -r .)
test: export GATEWAY_HOST=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_host | jq -r .)
test: export S3_URL=$(shell terraform -chdir=$(CHDIR) output -json s3_url | jq -r .)
test: export S3_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json s3_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=$(CHDIR) output -json dynamo_url | jq -r .)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json dynamo_token | jq -r .)
test:
	go test ./... -cover -race

load: export GATEWAY_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json gateway_token | jq -r .)
load: export GATEWAY_URL=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_url | jq -r .)
load: export GATEWAY_HOST=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_host | jq -r .)
load:
	k6 run k6/script.js

update:
	go get -u -t
	go mod tidy

docker:
	docker buildx build --platform linux/arm64 --platform linux/amd64 --tag ghcr.io/jhandguy/cloud-playground/gateway:latest --tag ghcr.io/jhandguy/cloud-playground/gateway:1.0.0 --file Dockerfile . --push
