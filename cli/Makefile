lint:
	golangci-lint run

build:
	go build .
	go test ./... -short -cover -race

test: export GATEWAY_API_KEY=$(shell terraform -chdir=../terraform output -json gateway_api_key | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=../terraform output -json gateway_url | jq -r .)
test:
	go test ./... -cover -race
	./cli message create -i "id" -c "content"
	./cli message get -i "id"
	./cli message delete -i "id"

update:
	go get -u
	go mod tidy