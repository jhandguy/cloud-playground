package object

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/cloud-playground/s3/opentelemetry"
	pb "github.com/jhandguy/cloud-playground/s3/pb/object"
)

type API struct {
	S3 S3
	grpc_health_v1.HealthServer
	pb.ObjectServiceServer
}

type S3 struct {
	Bucket string

	HeadBucketWithContext   func(ctx aws.Context, input *s3.HeadBucketInput, opts ...request.Option) (*s3.HeadBucketOutput, error)
	PutObjectWithContext    func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
	GetObjectWithContext    func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	DeleteObjectWithContext func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error)
}

func (api *API) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	switch req.GetService() {
	case "readiness":
		tracer := opentelemetry.GetTracer("object/CheckReadiness")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "object/CheckReadiness")
			defer span.End()
		}

		_, err := api.S3.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(api.S3.Bucket),
		})
		if err != nil {
			zap.S().Errorw("failed readiness check: %v", "error", err.Error())
			return &grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			}, nil
		}

		zap.S().Debugw("successfully checked readiness", "traceID", opentelemetry.GetTraceID(ctx))

		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	case "liveness":
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	default:
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}
}

func (api *API) CreateObject(ctx context.Context, req *pb.CreateObjectRequest) (*pb.CreateObjectResponse, error) {
	tracer := opentelemetry.GetTracer("object/CreateObject")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "object/CreateObject")
		defer span.End()
	}

	_, err := api.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(req.GetObject().GetId()),
		Body:   strings.NewReader(req.GetObject().GetContent()),
	})
	if err != nil {
		zap.S().Errorw("failed to create object", "error", err.Error())
		return nil, err
	}

	zap.S().Infow("successfully created object", "object", req.GetObject(), "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.CreateObjectResponse{Object: req.GetObject()}, nil
}

func (api *API) GetObject(ctx context.Context, req *pb.GetObjectRequest) (*pb.GetObjectResponse, error) {
	tracer := opentelemetry.GetTracer("object/GetObject")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "object/GetObject")
		defer span.End()
	}

	out, err := api.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(req.GetId()),
	})
	if err != nil {
		zap.S().Errorw("failed to get object", "error", err.Error())
		return nil, err
	}

	body := out.Body
	if body == nil {
		return &pb.GetObjectResponse{}, nil
	}

	defer func() {
		_ = body.Close()
	}()

	byt, err := io.ReadAll(body)
	if err != nil {
		zap.S().Errorw("failed to read object", "error", err.Error())
		return nil, err
	}

	object := &pb.Object{
		Id:      req.GetId(),
		Content: string(byt),
	}

	zap.S().Infow("successfully got object", "object", &object, "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.GetObjectResponse{
		Object: object,
	}, nil
}

func (api *API) DeleteObject(ctx context.Context, req *pb.DeleteObjectRequest) (*pb.DeleteObjectResponse, error) {
	tracer := opentelemetry.GetTracer("object/DeleteObject")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "object/DeleteObject")
		defer span.End()
	}

	object := &pb.Object{
		Id: req.GetId(),
	}

	_, err := api.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(object.GetId()),
	})
	if err != nil {
		zap.S().Errorw("failed to delete object", "error", err.Error())
		return nil, err
	}

	zap.S().Infow("successfully deleted object", "object", object, "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.DeleteObjectResponse{}, nil
}
