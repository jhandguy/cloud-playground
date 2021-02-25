lint:
	golangci-lint run

build:
	go build .
	go test ./... -short -cover -race

test: export GATEWAY_API_KEY=$(shell terraform -chdir=../terraform output -json gateway_api_key | jq -r .)
test: export GATEWAY_URL=$(shell terraform -chdir=../terraform output -json gateway_url | jq -r .)
test:
	go test ./... -cover -race
	./cli message create -u http://$$GATEWAY_URL -t $$GATEWAY_API_KEY -i "id" -c "content"
	./cli message get -u http://$$GATEWAY_URL -t $$GATEWAY_API_KEY -i "id"
	./cli message delete -u http://$$GATEWAY_URL -t $$GATEWAY_API_KEY -i "id"