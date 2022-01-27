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
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
	"github.com/jhandguy/devops-playground/gateway/opentelemetry"
	"github.com/jhandguy/devops-playground/gateway/prometheus"
)

func setupLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
}

func startTracing(ctx context.Context) {
	if endpoint := viper.GetString("tempo-url"); endpoint != "" {
		err := opentelemetry.StartTracing(ctx, endpoint)
		if err != nil {
			zap.S().Errorw("failed to start tracing", "error", err)
		}
	}
}

func stopTracing(ctx context.Context) {
	if err := opentelemetry.StopTracing(ctx); err != nil {
		zap.S().Errorw("failed to stop tracing", "error", err)
	}
}

func serveMetrics(path string) {
	port := viper.GetString("gateway-metrics-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen on metrics port", "error", err)
	}

	router := mux.NewRouter()
	router.Handle(path, promhttp.Handler()).Methods(http.MethodGet)
	if err := http.Serve(listener, router); err != nil {
		zap.S().Errorw("failed to serve metrics", "error", err)
	}
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

func newMessageAPI() *message.API {
	s3URL := viper.GetString("s3-url")
	s3Token := viper.GetString("s3-token")
	dynamoURL := viper.GetString("dynamo-url")
	dynamoToken := viper.GetString("dynamo-token")

	newS3Context := func(ctx context.Context) (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-token": s3Token})
		return context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 10*time.Second)
	}

	newS3ClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, s3URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	}

	newDynamoContext := func(ctx context.Context) (context.Context, context.CancelFunc) {
		md := metadata.New(map[string]string{"x-token": dynamoToken})
		return context.WithTimeout(metadata.NewOutgoingContext(ctx, md), 10*time.Second)
	}

	newDynamoClientConn := func(ctx context.Context) (*grpc.ClientConn, error) {
		return grpc.DialContext(ctx, dynamoURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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

func serveAPI(api *message.API, middlewares ...mux.MiddlewareFunc) {
	port := viper.GetString("gateway-http-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen on http port", "error", err)
	}

	router := routeAPI(api, middlewares...)
	if err := http.Serve(listener, router); err != nil {
		zap.S().Errorw("failed to serve API", "error", err)
	}
}

func main() {
	setupLogger()

	ctx := context.Background()
	startTracing(ctx)
	defer stopTracing(ctx)

	go serveMetrics("/monitoring/metrics")

	serveAPI(newMessageAPI(), prometheus.CollectMetrics, ensureValidToken)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
