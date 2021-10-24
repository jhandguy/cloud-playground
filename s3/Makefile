compile:
	protoc --go_out=. --go-grpc_out=. proto/object/object.pg.proto

lint:
	helm lint helm
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

GRPC_PORT ?= 8080
METRICS_PORT ?= 9090
test: export AWS_S3_ENDPOINT=$(shell terraform -chdir=../terraform output -json aws_s3_endpoint | jq -r .)
test: export AWS_S3_BUCKET=$(shell terraform -chdir=../terraform output -json aws_s3_bucket | jq -r .)
test: export S3_GRPC_PORT=$(GRPC_PORT)
test: export S3_METRICS_PORT=$(METRICS_PORT)
test: export S3_TOKEN=$(shell terraform -chdir=../terraform output -json s3_token | jq -r .)
test: export S3_URL=$(shell terraform -chdir=../terraform output -json s3_url | jq -r .)
test:
	go test ./... -cover -race

load: export S3_TOKEN=$(shell terraform -chdir=../terraform output -json s3_token | jq -r .)
load: export S3_URL=$(shell terraform -chdir=../terraform output -json s3_url | jq -r .)
load:
	k6 run k6/script.js

update:
	go get -u
	go mod tidy

docker:
	docker buildx build --tag ghcr.io/jhandguy/devops-playground/s3:latest --tag ghcr.io/jhandguy/devops-playground/s3:1.0.0 --file Dockerfile . --push
