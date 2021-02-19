package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
)

func retrieveEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("could not find environment variable %s", key)
	}
	return env
}

func newMessageAPI() *message.API {
	s3URL := retrieveEnv("S3_URL")
	s3Token := retrieveEnv("S3_TOKEN")
	dynamoURL := retrieveEnv("DYNAMO_URL")
	dynamoToken := retrieveEnv("DYNAMO_TOKEN")

	newS3Context := func() (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-api-key": s3Token})
		return context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 10*time.Second)
	}

	newS3ClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, s3URL, grpc.WithInsecure(), grpc.WithBlock())
	}

	newDynamoContext := func() (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-api-key": dynamoToken})
		return context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 10*time.Second)
	}

	newDynamoClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, dynamoURL, grpc.WithInsecure(), grpc.WithBlock())
	}

	return &message.API{
		ObjectAPI: &object.API{
			CreateObject: object.CreateObject(newS3Context, newS3ClientConn),
			GetObject:    object.GetObject(newS3Context, newS3ClientConn),
			DeleteObject: object.DeleteObject(newS3Context, newS3ClientConn),
		},
		ItemAPI: &item.API{
			CreateItem: item.CreateItem(newDynamoContext, newDynamoClientConn),
			GetItem:    item.GetItem(newDynamoContext, newDynamoClientConn),
			DeleteItem: item.DeleteItem(newDynamoContext, newDynamoClientConn),
		},
	}
}

func ensureValidAPIKey(next http.Handler) http.Handler {
	apiKey := retrieveEnv("GATEWAY_API_KEY")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == apiKey || strings.Contains(r.RequestURI, "/health") {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func serveAPI(api *message.API, middleware mux.MiddlewareFunc) *mux.Router {
	getHealth := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/health", getHealth).Methods(http.MethodGet)
	router.HandleFunc("/message", api.CreateMessage).Methods(http.MethodPost)
	router.HandleFunc("/message", api.GetMessage).Methods(http.MethodGet)
	router.HandleFunc("/message", api.DeleteMessage).Methods(http.MethodDelete)

	return router
}

func main() {
	port := retrieveEnv("GATEWAY_PORT")

	router := serveAPI(newMessageAPI(), ensureValidAPIKey)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
