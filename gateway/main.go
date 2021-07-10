package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
	"github.com/jhandguy/devops-playground/gateway/prometheus"
)

func newMessageAPI() *message.API {
	s3URL := viper.GetString("s3-url")
	s3Token := viper.GetString("s3-token")
	dynamoURL := viper.GetString("dynamo-url")
	dynamoToken := viper.GetString("dynamo-token")

	newS3Context := func() (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-token": s3Token})
		return context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 10*time.Second)
	}

	newS3ClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, s3URL, grpc.WithInsecure(), grpc.WithBlock())
	}

	newDynamoContext := func() (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-token": dynamoToken})
		return context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 10*time.Second)
	}

	newDynamoClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, dynamoURL, grpc.WithInsecure(), grpc.WithBlock())
	}

	return &message.API{
		ObjectAPI: &object.API{
			CheckHealth:  object.CheckHealth(newS3Context, newS3ClientConn),
			CreateObject: object.CreateObject(newS3Context, newS3ClientConn),
			GetObject:    object.GetObject(newS3Context, newS3ClientConn),
			DeleteObject: object.DeleteObject(newS3Context, newS3ClientConn),
		},
		ItemAPI: &item.API{
			CheckHealth: item.CheckHealth(newDynamoContext, newDynamoClientConn),
			CreateItem:  item.CreateItem(newDynamoContext, newDynamoClientConn),
			GetItem:     item.GetItem(newDynamoContext, newDynamoClientConn),
			DeleteItem:  item.DeleteItem(newDynamoContext, newDynamoClientConn),
		},
	}
}

func setDebugHeader(next http.Handler) http.Handler {
	deployment := viper.GetString("gateway-deployment")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-debug", deployment)
		next.ServeHTTP(w, r)
	})
}

func isValidToken(authorization, token string) bool {
	return strings.TrimPrefix(authorization, "Bearer ") == token
}

func ensureValidToken(next http.Handler) http.Handler {
	token := viper.GetString("gateway-token")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isValidToken(r.Header.Get("Authorization"), token) || strings.Contains(r.RequestURI, "/monitoring/") {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func serveMetrics(path string, listener net.Listener) {
	router := mux.NewRouter()
	router.Handle(path, promhttp.Handler()).Methods(http.MethodGet)

	if err := http.Serve(listener, router); err != nil {
		log.Fatalf("failed to serve metrics: %v", err)
	}
}

func routeAPI(api *message.API, middlewares ...mux.MiddlewareFunc) *mux.Router {
	setContentTypeHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}

	router := mux.NewRouter()
	router.Use(setContentTypeHeader)
	router.Use(middlewares...)
	router.HandleFunc("/monitoring/readiness", api.CheckReadiness).Methods(http.MethodGet)
	router.HandleFunc("/monitoring/liveness", api.CheckLiveness).Methods(http.MethodGet)
	router.HandleFunc("/message", api.CreateMessage).Methods(http.MethodPost)
	router.HandleFunc("/message/{id}", api.GetMessage).Methods(http.MethodGet)
	router.HandleFunc("/message/{id}", api.DeleteMessage).Methods(http.MethodDelete)

	return router
}

func serveAPI(api *message.API, listener net.Listener, middlewares ...mux.MiddlewareFunc) {
	router := routeAPI(api, middlewares...)

	if err := http.Serve(listener, router); err != nil {
		log.Fatalf("failed to serve API: %v", err)
	}
}

func main() {
	port := viper.GetString("gateway-metrics-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go serveMetrics("/monitoring/metrics", listener)

	port = viper.GetString("gateway-http-port")
	listener, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serveAPI(newMessageAPI(), listener, setDebugHeader, prometheus.CollectMetrics, ensureValidToken)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
