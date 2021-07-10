package item

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/devops-playground/gateway/pb/item"
)

type API struct {
	CheckHealth func(req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)
	CreateItem  func(req *item.CreateItemRequest) (*item.CreateItemResponse, error)
	GetItem     func(req *item.GetItemRequest) (*item.GetItemResponse, error)
	DeleteItem  func(req *item.DeleteItemRequest) (*item.DeleteItemResponse, error)
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

func CreateItem(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *item.CreateItemRequest) (*item.CreateItemResponse, error) {
	return func(req *item.CreateItemRequest) (*item.CreateItemResponse, error) {
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

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.CreateItem(ctx, req)
		if err != nil {
			log.Printf("failed to create item: %v", err)
			return nil, err
		}

		return resp, nil
	}
}

func GetItem(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *item.GetItemRequest) (*item.GetItemResponse, error) {
	return func(req *item.GetItemRequest) (*item.GetItemResponse, error) {
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

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.GetItem(ctx, req)
		if err != nil {
			log.Printf("failed to get item: %v", err)
			return nil, err
		}

		return resp, nil
	}
}

func DeleteItem(
	newContext func() (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(req *item.DeleteItemRequest) (*item.DeleteItemResponse, error) {
	return func(req *item.DeleteItemRequest) (*item.DeleteItemResponse, error) {
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

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.DeleteItem(ctx, req)
		if err != nil {
			log.Printf("failed to delete item: %v", err)
			return nil, err
		}

		return resp, nil
	}
}
