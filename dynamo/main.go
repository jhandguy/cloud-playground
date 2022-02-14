package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jhandguy/devops-playground/dynamo/item"
	"github.com/jhandguy/devops-playground/dynamo/opentelemetry"
	pb "github.com/jhandguy/devops-playground/dynamo/pb/item"
	"github.com/jhandguy/devops-playground/dynamo/prometheus"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
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
	port := viper.GetString("dynamo-metrics-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen on metrics port", "error", err)
	}

	http.Handle(path, promhttp.Handler())

	if err := http.Serve(listener, nil); err != nil {
		zap.S().Errorw("failed to serve metrics", "error", err)
	}
}

func isValidToken(authorization []string, token string) bool {
	if len(authorization) < 1 {
		return false
	}
	return strings.TrimPrefix(authorization[0], "Bearer ") == token
}

func ensureValidToken() grpc.UnaryServerInterceptor {
	token := viper.GetString("dynamo-token")

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		if !isValidToken(md["x-token"], token) {
			return nil, errInvalidToken
		}
		return handler(ctx, req)
	}
}

func newItemAPI() *item.API {
	endpoint := fmt.Sprintf("http://%s", viper.GetString("aws-dynamo-endpoint"))
	table := viper.GetString("aws-dynamo-table")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(endpoint),
		},
	}))

	client := dynamodb.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return &item.API{
		DynamoDB: item.DynamoDB{
			Table:                    table,
			DescribeTableWithContext: client.DescribeTableWithContext,
			PutItemWithContext:       client.PutItemWithContext,
			GetItemWithContext:       client.GetItemWithContext,
			DeleteItemWithContext:    client.DeleteItemWithContext,
		},
	}
}

func registerAPI(api *item.API, interceptors []grpc.UnaryServerInterceptor) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
	pb.RegisterItemServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, api)

	return s
}

func serveAPI(api *item.API, interceptors ...grpc.UnaryServerInterceptor) {
	port := viper.GetString("dynamo-grpc-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen on http port", "error", err)
	}

	s := registerAPI(api, interceptors)
	if err := s.Serve(listener); err != nil {
		zap.S().Errorw("failed to serve API", "error", err)
	}
}

func main() {
	setupLogger()

	ctx := context.Background()
	startTracing(ctx)
	defer stopTracing(ctx)

	go serveMetrics("/monitoring/metrics")

	serveAPI(newItemAPI(), prometheus.CollectMetrics, ensureValidToken())
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
