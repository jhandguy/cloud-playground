package main

import (
	"context"
	"dynamo/item"
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/test/bufconn"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
)

func TestIsValidToken(t *testing.T) {
	token := "token"

	auth := []string{
		fmt.Sprintf("Bearer %s", token),
	}
	assert.True(t, isValidToken(auth, token))

	auth = []string{}
	assert.False(t, isValidToken(auth, token))

	auth = []string{
		token,
	}
	assert.True(t, isValidToken(auth, token))

	auth = []string{
		"wrong",
	}
	assert.False(t, isValidToken(auth, token))
}

func TestServeAPI(t *testing.T) {
	var isPutItemWithContextCalled, isGetItemWithContextCalled, isDeleteItemWithContextCalled, isInterceptorCalled bool

	api := &item.API{
		DynamoDB: item.DynamoDB{
			Table: "table",
			PutItemWithContext: func(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
				isPutItemWithContextCalled = true
				return &dynamodb.PutItemOutput{}, nil
			},
			GetItemWithContext: func(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
				isGetItemWithContextCalled = true
				return &dynamodb.GetItemOutput{}, nil
			},
			DeleteItemWithContext: func(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
				isDeleteItemWithContextCalled = true
				return &dynamodb.DeleteItemOutput{}, nil
			},
		},
	}
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		isInterceptorCalled = true
		return handler(ctx, req)
	}

	bufSize := 1024 * 1024
	listener := bufconn.Listen(bufSize)

	go func() {
		serveAPI(api, interceptor, listener)
	}()

	ctx := context.Background()
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	c := item.NewItemServiceClient(conn)

	createReq := &item.CreateItemRequest{
		Name:    "name",
		Content: "content",
	}
	createRes, err := c.CreateItem(ctx, createReq)

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.True(t, isPutItemWithContextCalled)
	assert.True(t, isInterceptorCalled)

	getReq := &item.GetItemRequest{
		Id: "id",
	}
	getRes, err := c.GetItem(ctx, getReq)

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.True(t, isGetItemWithContextCalled)
	assert.True(t, isInterceptorCalled)

	deleteReq := &item.DeleteItemRequest{
		Id: "id",
	}
	deleteRes, err := c.DeleteItem(ctx, deleteReq)

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
	assert.True(t, isDeleteItemWithContextCalled)
	assert.True(t, isInterceptorCalled)
}

func TestService(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	go func() {
		main()
	}()

	host := retrieveEnv("DYNAMO_HOST")
	token := retrieveEnv("DYNAMO_TOKEN")

	md := metadata.New(map[string]string{"authorization": token})
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	c := item.NewItemServiceClient(conn)

	createReq := &item.CreateItemRequest{
		Name:    "name",
		Content: "content",
	}
	createRes, err := c.CreateItem(ctx, createReq)

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.NotNil(t, createRes.GetItem().GetId())
	assert.Equal(t, createRes.GetItem().GetName(), createReq.Name)
	assert.Equal(t, createRes.GetItem().GetContent(), createReq.Content)

	getReq := &item.GetItemRequest{
		Id: createRes.GetItem().GetId(),
	}
	getRes, err := c.GetItem(ctx, getReq)

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.Equal(t, getRes.GetItem().GetId(), getReq.GetId())
	assert.Equal(t, getRes.GetItem().GetName(), createRes.GetItem().GetName())
	assert.Equal(t, getRes.GetItem().GetContent(), createRes.GetItem().GetContent())

	deleteReq := &item.DeleteItemRequest{
		Id: getRes.GetItem().GetId(),
	}
	deleteRes, err := c.DeleteItem(ctx, deleteReq)

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
}
