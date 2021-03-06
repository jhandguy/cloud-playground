package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/jhandguy/devops-playground/dynamo/item"
	pb "github.com/jhandguy/devops-playground/dynamo/pb/item"
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
	interceptor := func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		isInterceptorCalled = true
		return handler(ctx, req)
	}

	bufSize := 1024 * 1024
	listener := bufconn.Listen(bufSize)

	go func() {
		serveAPI(api, listener, interceptor)
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

	c := pb.NewItemServiceClient(conn)

	createReq := &pb.CreateItemRequest{
		Item: &pb.Item{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}
	createRes, err := c.CreateItem(ctx, createReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.True(t, isPutItemWithContextCalled)
	assert.True(t, isInterceptorCalled)

	getReq := &pb.GetItemRequest{
		Id: createRes.GetItem().GetId(),
	}
	getRes, err := c.GetItem(ctx, getReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.True(t, isGetItemWithContextCalled)
	assert.True(t, isInterceptorCalled)

	deleteReq := &pb.DeleteItemRequest{
		Id: getRes.GetItem().GetId(),
	}
	deleteRes, err := c.DeleteItem(ctx, deleteReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
	assert.True(t, isDeleteItemWithContextCalled)
	assert.True(t, isInterceptorCalled)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	go main()

	port := viper.GetString("dynamo-grpc-port")
	url := fmt.Sprintf("localhost:%s", port)
	testDynamo(url, t)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := viper.GetString("dynamo-url")
	testDynamo(url, t)
}

func testDynamo(url string, t *testing.T) {
	token := viper.GetString("dynamo-token")

	md := metadata.New(map[string]string{"x-api-key": token})
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, url, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	c := pb.NewItemServiceClient(conn)

	createReq := &pb.CreateItemRequest{
		Item: &pb.Item{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}
	createRes, err := c.CreateItem(ctx, createReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.Equal(t, createReq.GetItem().GetId(), createRes.GetItem().GetId())
	assert.Equal(t, createReq.GetItem().GetContent(), createRes.GetItem().GetContent())

	getReq := &pb.GetItemRequest{
		Id: createRes.GetItem().GetId(),
	}
	getRes, err := c.GetItem(ctx, getReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.Equal(t, getReq.GetId(), getRes.GetItem().GetId())
	assert.Equal(t, createRes.GetItem().GetContent(), getRes.GetItem().GetContent())

	deleteReq := &pb.DeleteItemRequest{
		Id: getRes.GetItem().GetId(),
	}
	deleteRes, err := c.DeleteItem(ctx, deleteReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
}
