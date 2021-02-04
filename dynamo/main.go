package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jhandguy/devops-playground/dynamo/item"
	"github.com/jhandguy/devops-playground/dynamo/item/pb"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func retrieveEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("could not lookup env %s", key)
	}
	return env
}

func isValidToken(authorization []string, token string) bool {
	if len(authorization) < 1 {
		return false
	}
	return strings.TrimPrefix(authorization[0], "Bearer ") == token
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	token := retrieveEnv("DYNAMO_TOKEN")

	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !isValidToken(md["authorization"], token) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}

func newItemAPI() *item.API {
	endpoint := retrieveEnv("AWS_DYNAMO_ENDPOINT")
	table := retrieveEnv("AWS_DYNAMO_TABLE")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(endpoint),
		},
	}))

	client := dynamodb.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return &item.API{
		DynamoDB: item.DynamoDB{
			Table:                 table,
			PutItemWithContext:    client.PutItemWithContext,
			GetItemWithContext:    client.GetItemWithContext,
			DeleteItemWithContext: client.DeleteItemWithContext,
		},
	}
}

func serveAPI(api *item.API, interceptor grpc.UnaryServerInterceptor, listener net.Listener) {
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))

	pb.RegisterItemServiceServer(s, api)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	port := retrieveEnv("DYNAMO_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serveAPI(newItemAPI(), ensureValidToken, listener)
}
