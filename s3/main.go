package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"s3/object"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleObject(
	createObject func(http.ResponseWriter, *http.Request),
	getObject func(http.ResponseWriter, *http.Request),
	deleteObject func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createObject(w, r)
		case http.MethodGet:
			getObject(w, r)
		case http.MethodDelete:
			deleteObject(w, r)
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

func handleObjectFunc() func(http.ResponseWriter, *http.Request) {
	endpoint := retrieveEnv("AWS_S3_ENDPOINT")
	bucket := retrieveEnv("AWS_S3_BUCKET")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			Endpoint:         aws.String(endpoint),
		},
	}))

	client := s3.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return handleObject(
		object.CreateObject(object.CreateObjectFunc(client), bucket),
		object.GetObject(object.GetObjectFunc(client), bucket),
		object.DeleteObject(object.DeleteObjectFunc(client), bucket),
	)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	uriPrefix := retrieveEnv("URI_PREFIX")
	healthPath := retrieveEnv("HEALTH_PATH")

	http.HandleFunc(fmt.Sprintf("%s/object", uriPrefix), handleObjectFunc())
	http.HandleFunc(fmt.Sprintf("%s%s", uriPrefix, healthPath), checkHealth)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
