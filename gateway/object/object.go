package object

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/cloud-playground/gateway/opentelemetry"
	"github.com/jhandguy/cloud-playground/gateway/pb/object"
)

type API struct {
	CheckHealth  func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error)
	CreateObject func(ctx context.Context, req *object.CreateObjectRequest) (*object.CreateObjectResponse, error)
	GetObject    func(ctx context.Context, req *object.GetObjectRequest) (*object.GetObjectResponse, error)
	DeleteObject func(ctx context.Context, req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error)
}

func CheckHealth(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return func(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
		tracer := opentelemetry.GetTracer("message/CheckReadiness")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "object/CheckHealth")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		dynConn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err.Error())
			return nil, err
		}
		defer func() {
			_ = dynConn.Close()
		}()

		client := grpc_health_v1.NewHealthClient(dynConn)

		resp, err := client.Check(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to check health", "error", err.Error())
			return nil, err
		}

		return resp, nil
	}
}

func CreateObject(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *object.CreateObjectRequest) (*object.CreateObjectResponse, error) {
	return func(ctx context.Context, req *object.CreateObjectRequest) (*object.CreateObjectResponse, error) {
		tracer := opentelemetry.GetTracer("message/CreateMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "object/CreateObject")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err.Error())
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.CreateObject(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to create object", "error", err.Error())
			return nil, err
		}

		zap.S().Infow("successfully created object", "object", resp.GetObject(), "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}

func GetObject(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *object.GetObjectRequest) (*object.GetObjectResponse, error) {
	return func(ctx context.Context, req *object.GetObjectRequest) (*object.GetObjectResponse, error) {
		tracer := opentelemetry.GetTracer("message/GetMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "object/GetObject")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err.Error())
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.GetObject(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to get object", "error", err.Error())
			return nil, err
		}

		zap.S().Infow("successfully got object", "object", resp.GetObject(), "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}

func DeleteObject(
	newContext func(ctx context.Context) (context.Context, context.CancelFunc),
	newClientConn func(ctx context.Context) (*grpc.ClientConn, error),
) func(ctx context.Context, req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error) {
	return func(ctx context.Context, req *object.DeleteObjectRequest) (*object.DeleteObjectResponse, error) {
		tracer := opentelemetry.GetTracer("message/DeleteMessage")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "object/DeleteObject")
			defer span.End()
		}

		ctx, cancel := newContext(ctx)
		defer cancel()

		s3Conn, err := newClientConn(ctx)
		if err != nil {
			zap.S().Errorw("failed to dial", "error", err.Error())
			return nil, err
		}
		defer func() {
			_ = s3Conn.Close()
		}()

		client := object.NewObjectServiceClient(s3Conn)

		resp, err := client.DeleteObject(ctx, req)
		if err != nil {
			zap.S().Errorw("failed to delete object", "error", err.Error())
			return nil, err
		}

		zap.S().Infow("successfully deleted object", "object", object.Object{Id: req.GetId()}, "traceID", opentelemetry.GetTraceID(ctx))

		return resp, nil
	}
}
