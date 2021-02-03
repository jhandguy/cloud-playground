package item

import (
	"context"
	"dynamo/item/pb"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestCreateItem(t *testing.T) {
	var actTable, actId, actName, actContent string

	api := API{
		DynamoDB: DynamoDB{
			Table: "table",
			PutItemWithContext: func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
				actTable = *input.TableName
				actId = *input.Item["id"].S
				actName = *input.Item["name"].S
				actContent = *input.Item["content"].S

				return &dynamodb.PutItemOutput{}, nil
			},
			GetItemWithContext:    nil,
			DeleteItemWithContext: nil,
		},
	}

	req := &pb.CreateItemRequest{
		Name:    "name",
		Content: "content",
	}
	resp, err := api.CreateItem(context.Background(), req)

	assert.Equal(t, actTable, api.DynamoDB.Table)
	assert.Equal(t, actName, req.Name)
	assert.Equal(t, actContent, req.Content)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Item.Id, actId)
	assert.Equal(t, resp.Item.Name, req.Name)
	assert.Equal(t, resp.Item.Content, req.Content)
}

func TestGetItem(t *testing.T) {
	var actTable, actId string
	item := pb.Item{
		Id:      "id",
		Name:    "name",
		Content: "content",
	}

	api := API{
		DynamoDB: DynamoDB{
			Table:              "table",
			PutItemWithContext: nil,
			GetItemWithContext: func(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
				actTable = *input.TableName
				actId = *input.Key["id"].S

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
		Id: item.Id,
	}
	resp, err := api.GetItem(context.Background(), req)

	assert.Equal(t, actTable, api.DynamoDB.Table)
	assert.Equal(t, actId, req.Id)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Item.Id, item.Id)
	assert.Equal(t, resp.Item.Name, item.Name)
	assert.Equal(t, resp.Item.Content, item.Content)
}

func TestDeleteItem(t *testing.T) {
	var actTable, actId string

	api := API{
		DynamoDB: DynamoDB{
			Table:              "table",
			PutItemWithContext: nil,
			GetItemWithContext: nil,
			DeleteItemWithContext: func(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
				actTable = *input.TableName
				actId = *input.Key["id"].S

				return &dynamodb.DeleteItemOutput{}, nil
			},
		},
	}

	req := &pb.DeleteItemRequest{
		Id: "id",
	}
	resp, err := api.DeleteItem(context.Background(), req)

	assert.Equal(t, actTable, api.DynamoDB.Table)
	assert.Equal(t, actId, req.Id)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
