package item

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/devops-playground/gateway/opentelemetry"
	"github.com/jhandguy/devops-playground/gateway/pb/item"
)

type API struct {
	CheckHealth func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)
	CreateItem  func(ctx context.Context, req *item.CreateItemRequest) (*item.CreateItemResponse, error)
	GetItem     func(ctx context.Context, req *item.GetItemRequest) (*item.GetItemResponse, error)
	DeleteItem  func(ctx context.Context, req *item.DeleteItemRequest) (*item.DeleteItemResponse, error)
}

func CheckHealth(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
		tracer := opentelemetry.GetTracer("message/CheckReadiness")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "item/CheckHealth")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err)
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := grpc_health_v1.NewHealthClient(dynConn)

		resp, err := client.Check(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to check health", "error", err)
			return nil, err
		}

		return resp, nil
	}
}

func CreateItem(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *item.CreateItemRequest) (*item.CreateItemResponse, error) {
	return func(ctx context.Context, req *item.CreateItemRequest) (*item.CreateItemResponse, error) {
		tracer := opentelemetry.GetTracer("message/CreateMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "item/CreateItem")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err)
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.CreateItem(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to create item", "error", err)
			return nil, err
		}

		zap.S().Infow("successfully created item", "item", resp.GetItem(), "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}

func GetItem(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *item.GetItemRequest) (*item.GetItemResponse, error) {
	return func(ctx context.Context, req *item.GetItemRequest) (*item.GetItemResponse, error) {
		tracer := opentelemetry.GetTracer("message/GetMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "item/GetItem")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err)
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.GetItem(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to get item", "error", err)
			return nil, err
		}

		zap.S().Infow("successfully got item", "item", resp.GetItem(), "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}

func DeleteItem(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *item.DeleteItemRequest) (*item.DeleteItemResponse, error) {
	return func(ctx context.Context, req *item.DeleteItemRequest) (*item.DeleteItemResponse, error) {
		tracer := opentelemetry.GetTracer("message/DeleteMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "item/DeleteItem")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err)
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := item.NewItemServiceClient(dynConn)

		resp, err := client.DeleteItem(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to delete item", "error", err)
			return nil, err
		}

		zap.S().Infow("successfully deleted item", "item", item.Item{Id: req.GetId()}, "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}
