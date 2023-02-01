ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
GRPC_PORT ?= 8080
METRICS_PORT ?= 9090

compile:
	protoc --go_out=. --go-grpc_out=. proto/object/object.pg.proto

lint_helm:
	helm lint helm

lint:
	golangci-lint run

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
	go test ./... -short -cover -race

test: export AWS_S3_ENDPOINT=$(shell terraform -chdir=$(CHDIR) output -json aws_s3_endpoint | jq -r .)
test: export AWS_S3_BUCKET=$(shell terraform -chdir=$(CHDIR) output -json aws_s3_bucket | jq -r .)
test: export S3_GRPC_PORT=$(GRPC_PORT)
test: export S3_METRICS_PORT=$(METRICS_PORT)
test: export S3_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json s3_token | jq -r .)
test: export S3_URL=$(shell terraform -chdir=$(CHDIR) output -json s3_url | jq -r .)
test:
	go test ./... -cover -race

load: export S3_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json s3_token | jq -r .)
load: export S3_URL=$(shell terraform -chdir=$(CHDIR) output -json s3_url | jq -r .)
load:
	k6 run k6/script.js

update:
	go get -u -t
	go mod tidy

docker:
	docker buildx build --platform linux/arm64 --platform linux/amd64 --tag ghcr.io/jhandguy/cloud-playground/s3:latest --tag ghcr.io/jhandguy/cloud-playground/s3:1.0.0 --file Dockerfile . --push
