package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
)

func newMessageAPI() *message.API {
	s3URL := viper.GetString("s3-url")
	s3Token := viper.GetString("s3-token")
	dynamoURL := viper.GetString("dynamo-url")
	dynamoToken := viper.GetString("dynamo-token")

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

func isValidAPIKey(authorization, apiKey string) bool {
	return strings.TrimPrefix(authorization, "Bearer ") == apiKey
}

func ensureValidAPIKey(next http.Handler) http.Handler {
	apiKey := viper.GetString("gateway-api-key")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isValidAPIKey(r.Header.Get("Authorization"), apiKey) || strings.Contains(r.RequestURI, "/health") {
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

	setContentTypeHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}

	router := mux.NewRouter()
	router.Use(setContentTypeHeader)
	router.Use(middleware)
	router.HandleFunc("/health", getHealth).Methods(http.MethodGet)
	router.HandleFunc("/message", api.CreateMessage).Methods(http.MethodPost)
	router.HandleFunc("/message/{id}", api.GetMessage).Methods(http.MethodGet)
	router.HandleFunc("/message/{id}", api.DeleteMessage).Methods(http.MethodDelete)

	return router
}

func main() {
	port := viper.GetString("gateway-port")

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

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
