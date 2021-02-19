package item

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	pb "github.com/jhandguy/devops-playground/dynamo/pb/item"
)

type API struct {
	DynamoDB DynamoDB
	pb.ItemServiceServer
}

type DynamoDB struct {
	Table string

	PutItemWithContext    func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	GetItemWithContext    func(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)
	DeleteItemWithContext func(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error)
}

func (api *API) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	item := &pb.Item{
		Id:      uuid.New().String(),
		Name:    req.Name,
		Content: req.Content,
	}

	it, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Printf("failed to marshal item: %v", err)
		return nil, err
	}

	_, err = api.DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      it,
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		log.Printf("failed to create item: %v", err)
		return nil, err
	}

	return &pb.CreateItemResponse{
		Item: item,
	}, nil
}

func (api *API) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	out, err := api.DynamoDB.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(req.Id),
			},
		},
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		log.Printf("failed to get item: %v", err)
		return nil, err
	}

	var item pb.Item
	err = dynamodbattribute.UnmarshalMap(out.Item, &item)
	if err != nil {
		log.Printf("failed to unmarshal item: %v", err)
		return nil, err
	}

	return &pb.GetItemResponse{
		Item: &item,
	}, nil
}

func (api *API) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	_, err := api.DynamoDB.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(req.Id),
			},
		},
		TableName: aws.String(api.DynamoDB.Table),
	})
	if err != nil {
		log.Printf("failed to delete item: %v", err)
		return nil, err
	}

	return &pb.DeleteItemResponse{}, nil
}
