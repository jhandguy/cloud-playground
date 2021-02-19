package item

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	pb "github.com/jhandguy/devops-playground/dynamo/pb/item"
)

func TestCreateItem(t *testing.T) {
	var actTable, actID, actName, actContent string

	api := API{
		DynamoDB: DynamoDB{
			Table: "table",
			PutItemWithContext: func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
				actTable = *input.TableName
				actID = *input.Item["id"].S
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

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.Name, actName)
	assert.Equal(t, req.Content, actContent)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, actID, resp.Item.Id)
	assert.Equal(t, req.Name, resp.Item.Name)
	assert.Equal(t, req.Content, resp.Item.Content)
}

func TestGetItem(t *testing.T) {
	var actTable, actID string
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
		Id: item.Id,
	}
	resp, err := api.GetItem(context.Background(), req)

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.Id, actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, item.Id, resp.Item.Id)
	assert.Equal(t, item.Name, resp.Item.Name)
	assert.Equal(t, item.Content, resp.Item.Content)
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
		Id: "id",
	}
	resp, err := api.DeleteItem(context.Background(), req)

	assert.Equal(t, api.DynamoDB.Table, actTable)
	assert.Equal(t, req.Id, actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
