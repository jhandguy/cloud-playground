package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jhandguy/devops-playground/s3/object"
	pb "github.com/jhandguy/devops-playground/s3/pb/object"
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

	if !isValidToken(md["x-api-key"], token) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}

func newObjectAPI() *object.API {
	endpoint := viper.GetString("aws-s3-endpoint")
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
			PutObjectWithContext:    client.PutObjectWithContext,
			GetObjectWithContext:    client.GetObjectWithContext,
			DeleteObjectWithContext: client.DeleteObjectWithContext,
		},
	}
}

func serveAPI(api *object.API, interceptor grpc.UnaryServerInterceptor, listener net.Listener) {
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))

	pb.RegisterObjectServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	port := viper.GetString("s3-port")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serveAPI(newObjectAPI(), ensureValidToken, listener)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
