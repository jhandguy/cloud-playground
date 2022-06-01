package item

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jhandguy/cloud-playground/dynamo/opentelemetry"
	pb "github.com/jhandguy/cloud-playground/dynamo/pb/item"
)

type API struct {
	DynamoDB DynamoDB
	grpc_health_v1.HealthServer
	pb.ItemServiceServer
}

type DynamoDB struct {
	Table string

	DescribeTableWithContext func(ctx aws.Context, input *dynamodb.DescribeTableInput, opts ...request.Option) (*dynamodb.DescribeTableOutput, error)
	PutItemWithContext       func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	GetItemWithContext       func(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)
	DeleteItemWithContext    func(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error)
}

func (api *API) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	switch req.GetService() {
	case "readiness":
		tracer := opentelemetry.GetTracer("item/CheckReadiness")
		if tracer != nil {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "item/CheckReadiness")
			defer span.End()
		}

		_, err := api.DynamoDB.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(api.DynamoDB.Table),
		})
		if err != nil {
			zap.S().Errorw("failed readiness check", "error", err.Error())
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

func (api *API) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	tracer := opentelemetry.GetTracer("item/CreateItem")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "item/CreateItem")
		defer span.End()
	}

	it, err := dynamodbattribute.MarshalMap(req.GetItem())
	if err != nil {
		zap.S().Errorw("failed to marshal item", "error", err.Error())
		return nil, err
	}

	_, err = api.DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      it,
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		zap.S().Errorw("failed to create item", "error", err.Error())
		return nil, err
	}

	zap.S().Infow("successfully created item", "item", req.GetItem(), "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.CreateItemResponse{
		Item: req.GetItem(),
	}, nil
}

func (api *API) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	tracer := opentelemetry.GetTracer("item/GetItem")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "item/GetItem")
		defer span.End()
	}

	out, err := api.DynamoDB.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(req.GetId()),
			},
		},
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		zap.S().Errorw("failed to get item", "error", err.Error())
		return nil, err
	}

	var item pb.Item
	err = dynamodbattribute.UnmarshalMap(out.Item, &item)
	if err != nil {
		zap.S().Errorw("failed to unmarshal item", "error", err.Error())
		return nil, err
	}

	zap.S().Infow("successfully got item", "item", &item, "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.GetItemResponse{
		Item: &item,
	}, nil
}

func (api *API) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	tracer := opentelemetry.GetTracer("item/DeleteItem")
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "item/DeleteItem")
		defer span.End()
	}

	item := &pb.Item{
		Id: req.GetId(),
	}

	_, err := api.DynamoDB.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(item.GetId()),
			},
		},
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		zap.S().Errorw("failed to delete item", "error", err.Error())
		return nil, err
	}

	zap.S().Infow("successfully deleted item", "item", item, "traceID", opentelemetry.GetTraceID(ctx))

	return &pb.DeleteItemResponse{}, nil
}
