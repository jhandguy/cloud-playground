ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
GRPC_PORT ?= 8080
METRICS_PORT ?= 9090

compile:
	protoc --go_out=. --go-grpc_out=. proto/item/item.pg.proto

lint:
	helm lint helm
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

test: export AWS_DYNAMO_ENDPOINT=$(shell terraform -chdir=$(CHDIR) output -json aws_dynamo_endpoint | jq -r .)
test: export AWS_DYNAMO_TABLE=$(shell terraform -chdir=$(CHDIR) output -json aws_dynamo_table | jq -r .)
test: export DYNAMO_GRPC_PORT=$(GRPC_PORT)
test: export DYNAMO_METRICS_PORT=$(METRICS_PORT)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json dynamo_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=$(CHDIR) output -json dynamo_url | jq -r .)
test:
	go test ./... -cover -race

load: export DYNAMO_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json dynamo_token | jq -r .)
load: export DYNAMO_URL=$(shell terraform -chdir=$(CHDIR) output -json dynamo_url | jq -r .)
load:
	k6 run k6/script.js

update:
	go get -u -t
	go mod tidy

docker:
	docker buildx build --tag ghcr.io/jhandguy/devops-playground/dynamo:latest --tag ghcr.io/jhandguy/devops-playground/dynamo:1.0.0 --file Dockerfile . --push
