ENVIRONMENT ?=
CHDIR = ../terraform/environments/$(ENVIRONMENT)
ROUNDS ?= 100

lint:
	helm lint helm
	golangci-lint run

build:
	$(shell if [ $$(uname) != "Linux" ]; then echo "env GOOS=linux GOARCH=amd64"; fi) go build .
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
	docker buildx build --tag ghcr.io/jhandguy/devops-playground/cli:latest --tag ghcr.io/jhandguy/devops-playground/cli:1.0.0 --file Dockerfile . --push
