package item

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pb "github.com/jhandguy/cloud-playground/dynamo/pb/item"
)

func TestCreateItem(t *testing.T) {
	var actTable, actID, actContent string

	api := API{
		DynamoDB: DynamoDB{
			Table: "table",
			PutItemWithContext: func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
				actTable = *input.TableName
				actID = *input.Item["id"].S
				actContent = *input.Item["content"].S

				return &dynamodb.PutItemOutput{}, nil
			},
			GetItemWithContext:    nil,
			DeleteItemWithContext: nil,
		},
	}

	req := &pb.CreateItemRequest{
		Item: &pb.Item{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}
	resp, err := api.CreateItem(context.Background(), req)

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.GetItem().GetId(), actID)
	assert.Equal(t, req.GetItem().GetContent(), actContent)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.GetItem().GetId(), resp.Item.GetId())
	assert.Equal(t, req.GetItem().GetContent(), resp.Item.GetContent())
}

func TestGetItem(t *testing.T) {
	var actTable, actID string
	item := pb.Item{
		Id:      uuid.NewString(),
		Content: "content",
	}

	api := API{
		DynamoDB: DynamoDB{
			Table:              "table",
			PutItemWithContext: nil,
			GetItemWithContext: func(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
				actTable = *input.TableName
				actID = *input.Key["id"].S

				itm, err := dynamodbattribute.MarshalMap(&item)
				if err != nil {
					return nil, err
				}

				return &dynamodb.GetItemOutput{
					Item: itm,
				}, nil
			},
			DeleteItemWithContext: nil,
		},
	}

	req := &pb.GetItemRequest{
		Id: item.GetId(),
	}
	resp, err := api.GetItem(context.Background(), req)

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.GetId(), actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, item.GetId(), resp.Item.GetId())
	assert.Equal(t, item.GetContent(), resp.Item.GetContent())
}

func TestDeleteItem(t *testing.T) {
	var actTable, actID string

	api := API{
		DynamoDB: DynamoDB{
			Table:              "table",
			PutItemWithContext: nil,
			GetItemWithContext: nil,
			DeleteItemWithContext: func(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
				actTable = *input.TableName
				actID = *input.Key["id"].S

				return &dynamodb.DeleteItemOutput{}, nil
			},
		},
	}

	req := &pb.DeleteItemRequest{
		Id: uuid.NewString(),
	}
	resp, err := api.DeleteItem(context.Background(), req)

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.GetId(), actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
