ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
ROUNDS ?= 100

lint_helm:
	helm lint helm

lint:
	golangci-lint run

format:
	go fmt ./...

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
	go test ./... -short -cover -race

test: export GATEWAY_TOKEN=$(shell terraform -chdir=$(CHDIR) output -json gateway_token | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_url | jq -r .)
test: export GATEWAY_HOST=$(shell terraform -chdir=$(CHDIR) output -json ingress_gateway_host | jq -r .)
test: export PUSHGATEWAY_URL=$(shell terraform -chdir=$(CHDIR) output -json pushgateway_url | jq -r .)
test:
	go test ./... -cover -race
	go run main.go message create -i "id" -c "content"
	go run main.go message get -i "id"
	go run main.go message delete -i "id"
	go run main.go load test -r $(ROUNDS)

update:
	go get -u -t
	go mod tidy

docker:
	docker buildx build --platform linux/arm64 --platform linux/amd64 --tag ghcr.io/jhandguy/cloud-playground/cli:latest --tag ghcr.io/jhandguy/cloud-playground/cli:1.0.0 --file Dockerfile . --push
