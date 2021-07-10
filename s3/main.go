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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jhandguy/devops-playground/s3/object"
	pb "github.com/jhandguy/devops-playground/s3/pb/object"
	"github.com/jhandguy/devops-playground/s3/prometheus"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

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

func serveMetrics(path string, listener net.Listener) {
	http.Handle(path, promhttp.Handler())

	if err := http.Serve(listener, nil); err != nil {
		log.Fatalf("failed to serve metrics: %v", err)
	}
}

func serveAPI(api *object.API, listener net.Listener, interceptors ...grpc.UnaryServerInterceptor) {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	pb.RegisterObjectServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, api)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve API: %v", err)
	}
}

func main() {
	port := viper.GetString("s3-metrics-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go serveMetrics("/monitoring/metrics", listener)

	port = viper.GetString("s3-grpc-port")
	listener, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serveAPI(newObjectAPI(), listener, prometheus.CollectMetrics, ensureValidToken)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
