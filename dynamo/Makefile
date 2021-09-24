compile:
	protoc --go_out=. --go-grpc_out=. proto/item/item.pg.proto

lint:
	helm lint helm
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

GRPC_PORT ?= 8080
METRICS_PORT ?= 9090
test: export AWS_DYNAMO_ENDPOINT=$(shell terraform -chdir=../terraform output -json aws_dynamo_endpoint | jq -r .)
test: export AWS_DYNAMO_TABLE=$(shell terraform -chdir=../terraform output -json aws_dynamo_table | jq -r .)
test: export DYNAMO_GRPC_PORT=$(GRPC_PORT)
test: export DYNAMO_METRICS_PORT=$(METRICS_PORT)
test: export DYNAMO_TOKEN=$(shell terraform -chdir=../terraform output -json dynamo_token | jq -r .)
test: export DYNAMO_URL=$(shell terraform -chdir=../terraform output -json dynamo_url | jq -r .)
test:
	go test ./... -cover -race

update:
	go get -u
	go mod tidy

docker:
	docker buildx build --tag ghcr.io/jhandguy/devops-playground/dynamo:latest --tag ghcr.io/jhandguy/devops-playground/dynamo:1.0.0 --file Dockerfile . --push
