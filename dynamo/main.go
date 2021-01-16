package main

import (
	"dynamo/item"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func handleItem(
	createItem func(http.ResponseWriter, *http.Request),
	getItem func(http.ResponseWriter, *http.Request),
	deleteItem func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createItem(w, r)
		case http.MethodGet:
			getItem(w, r)
		case http.MethodDelete:
			deleteItem(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func retrieveEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("could not find environment variable %s", key)
	}
	return env
}

func handleItemFunc() func(http.ResponseWriter, *http.Request) {
	endpoint := retrieveEnv("AWS_DYNAMO_ENDPOINT")
	table := retrieveEnv("AWS_DYNAMO_TABLE")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(endpoint),
		},
	}))

	client := dynamodb.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return handleItem(
		item.CreateItem(item.CreateItemFunc(client), table),
		item.GetItem(item.GetItemFunc(client), table),
		item.DeleteItem(item.DeleteItemFunc(client), table),
	)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func authMiddleware(apiKey string, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == apiKey {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

func main() {
	uriPrefix := retrieveEnv("URI_PREFIX")
	healthPath := retrieveEnv("HEALTH_PATH")
	apiKey := retrieveEnv("API_KEY")

	http.HandleFunc(fmt.Sprintf("%s/item", uriPrefix), authMiddleware(apiKey, handleItemFunc()))
	http.HandleFunc(fmt.Sprintf("%s%s", uriPrefix, healthPath), checkHealth)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
