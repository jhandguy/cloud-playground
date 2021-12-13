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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jhandguy/devops-playground/s3/object"
	"github.com/jhandguy/devops-playground/s3/opentelemetry"
	pb "github.com/jhandguy/devops-playground/s3/pb/object"
	"github.com/jhandguy/devops-playground/s3/prometheus"
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
	port := viper.GetString("s3-metrics-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen", "error", err)
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

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	token := viper.GetString("s3-token")

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

func newObjectAPI() *object.API {
	endpoint := fmt.Sprintf("http://%s", viper.GetString("aws-s3-endpoint"))
	bucket := viper.GetString("aws-s3-bucket")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			Endpoint:         aws.String(endpoint),
		},
	}))

	client := s3.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return &object.API{
		S3: object.S3{
			Bucket:                  bucket,
			HeadBucketWithContext:   client.HeadBucketWithContext,
			PutObjectWithContext:    client.PutObjectWithContext,
			GetObjectWithContext:    client.GetObjectWithContext,
			DeleteObjectWithContext: client.DeleteObjectWithContext,
		},
	}
}

func registerAPI(api *object.API, interceptors []grpc.UnaryServerInterceptor) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
	pb.RegisterObjectServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, api)

	return s
}

func serveAPI(api *object.API, interceptors ...grpc.UnaryServerInterceptor) {
	port := viper.GetString("s3-grpc-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zap.S().Errorw("failed to listen on http port", "error", err)
	}

	s := registerAPI(api, interceptors)
	if err := s.Serve(listener); err != nil {
		zap.S().Errorw("failed to serve API", err)
	}
}

func main() {
	setupLogger()

	ctx := context.Background()
	startTracing(ctx)
	defer stopTracing(ctx)

	go serveMetrics("/monitoring/metrics")

	serveAPI(newObjectAPI(), prometheus.CollectMetrics, ensureValidToken)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
