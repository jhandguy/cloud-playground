compile:
	protoc --go_out=. --go-grpc_out=. proto/object/object.pg.proto

lint:
	golangci-lint run

build:
	$(shell if [[ $OSTYPE != "linux-gnu"* ]]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
	go test ./... -short -cover -race

test: export AWS_S3_ENDPOINT=$(shell terraform -chdir=../terraform output -json aws_s3_endpoint | jq -r .)
test: export AWS_S3_BUCKET=$(shell terraform -chdir=../terraform output -json aws_s3_bucket | jq -r .)
test: export S3_HOST=$(shell terraform -chdir=../terraform output -json s3_host | jq -r .)
test: export S3_PORT=$(shell terraform -chdir=../terraform output -json s3_port | jq -r .)
test: export S3_TOKEN=$(shell terraform -chdir=../terraform output -json s3_token | jq -r .)
test:
	go test ./... -cover -race