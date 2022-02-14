package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	handler := promhttp.Handler()
	getMetrics := func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}

	router := gin.New()
	router.GET(path, getMetrics)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		zap.S().Errorw("failed to serve metrics", "error", err)
	}
}

func isValidToken(authorization, token string) bool {
	return strings.TrimPrefix(authorization, "Bearer ") == token
}

func ensureValidToken() gin.HandlerFunc {
	token := viper.GetString("gateway-token")
	return func(c *gin.Context) {
		if isValidToken(c.Request.Header.Get("Authorization"), token) || strings.Contains(c.Request.RequestURI, "/monitoring/") {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		}
	}
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

func routeAPI(api *message.API, middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/monitoring/readiness", api.CheckReadiness)
	router.GET("/monitoring/liveness", api.CheckLiveness)
	router.POST("/message", api.CreateMessage)
	router.GET("/message/:id", api.GetMessage)
	router.DELETE("/message/:id", api.DeleteMessage)

	return router
}

func serveAPI(api *message.API, middlewares ...gin.HandlerFunc) {
	port := viper.GetString("gateway-http-port")
	router := routeAPI(api, middlewares...)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		zap.S().Errorw("failed to serve API", "error", err)
	}
}

func main() {
	setupLogger()

	ctx := context.Background()
	startTracing(ctx)
	defer stopTracing(ctx)

	go serveMetrics("/monitoring/metrics")

	serveAPI(newMessageAPI(), prometheus.CollectMetrics, ensureValidToken())
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	gin.SetMode(gin.ReleaseMode)
}
