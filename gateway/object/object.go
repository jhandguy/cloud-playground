package object

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/devops-playground/gateway/pb/object"
)

type API struct {
	CheckHealth  func(req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)
	CreateObject func(req *object.CreateObjectRequest) (*object.CreateObjectResponse, error)
	GetObject    func(req *object.GetObjectRequest) (*object.GetObjectResponse, error)
	DeleteObject func(req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error)
}

func CheckHealth(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return func(req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
		ctx, cancel := newContext()
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			log.Printf("failed to dial: %v", err)
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := grpc_health_v1.NewHealthClient(dynConn)

		resp, err := client.Check(ctx, req)
		if err != nil {
			log.Printf("failed health check: %v", err)
			return nil, err
		}

		return resp, nil
	}
}

func CreateObject(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *object.CreateObjectRequest) (*object.CreateObjectResponse, error) {
	return func(req *object.CreateObjectRequest) (*object.CreateObjectResponse, error) {
		ctx, cancel := newContext()
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			log.Printf("failed to dial: %v", err)
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.CreateObject(ctx, req)
		if err != nil {
			log.Printf("failed to create object: %v", err)
			return nil, err
		}

		return resp, nil
	}
}

func GetObject(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *object.GetObjectRequest) (*object.GetObjectResponse, error) {
	return func(req *object.GetObjectRequest) (*object.GetObjectResponse, error) {
		ctx, cancel := newContext()
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			log.Printf("failed to dial: %v", err)
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.GetObject(ctx, req)
		if err != nil {
			log.Printf("failed to get object: %v", err)
			return nil, err
		}

		return resp, nil
	}
}

func DeleteObject(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error) {
	return func(req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error) {
		ctx, cancel := newContext()
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			log.Printf("failed to dial: %v", err)
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.DeleteObject(ctx, req)
		if err != nil {
			log.Printf("failed to delete object: %v", err)
			return nil, err
		}

		return resp, nil
	}
}
