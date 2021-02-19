package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/jhandguy/devops-playground/s3/object"
	pb "github.com/jhandguy/devops-playground/s3/pb/object"
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
	var isPutObjectWithContextCalled, isGetObjectWithContextCalled, isDeleteObjectWithContextCalled, isInterceptorCalled bool

	api := &object.API{
		S3: object.S3{
			Bucket: "bucket",
			PutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				isPutObjectWithContextCalled = true
				return &s3.PutObjectOutput{}, nil
			},
			GetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				isGetObjectWithContextCalled = true
				return &s3.GetObjectOutput{}, nil
			},
			DeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				isDeleteObjectWithContextCalled = true
				return &s3.DeleteObjectOutput{}, nil
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

	c := pb.NewObjectServiceClient(conn)

	createReq := &pb.CreateObjectRequest{
		Name:    "name",
		Content: "content",
	}
	createRes, err := c.CreateObject(ctx, createReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.True(t, isPutObjectWithContextCalled)
	assert.True(t, isInterceptorCalled)

	getReq := &pb.GetObjectRequest{
		Name: "name",
	}
	getRes, err := c.GetObject(ctx, getReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.True(t, isGetObjectWithContextCalled)
	assert.True(t, isInterceptorCalled)

	deleteReq := &pb.DeleteObjectRequest{
		Name: "name",
	}
	deleteRes, err := c.DeleteObject(ctx, deleteReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
	assert.True(t, isDeleteObjectWithContextCalled)
	assert.True(t, isInterceptorCalled)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	go main()

	port := retrieveEnv("S3_PORT")
	url := fmt.Sprintf("localhost:%s", port)
	testS3(url, t)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := retrieveEnv("S3_URL")
	testS3(url, t)
}

func testS3(url string, t *testing.T) {
	token := retrieveEnv("S3_TOKEN")

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

	c := pb.NewObjectServiceClient(conn)

	createReq := &pb.CreateObjectRequest{
		Name:    "name",
		Content: "content",
	}
	createRes, err := c.CreateObject(ctx, createReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, createRes)
	assert.Equal(t, createReq.Name, createRes.GetObject().GetName())
	assert.Equal(t, createReq.Content, createRes.GetObject().GetContent())

	getReq := &pb.GetObjectRequest{
		Name: createReq.GetName(),
	}
	getRes, err := c.GetObject(ctx, getReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, getRes)
	assert.Equal(t, createRes.GetObject().GetName(), getRes.GetObject().GetName())
	assert.Equal(t, createRes.GetObject().GetContent(), getRes.GetObject().GetContent())

	deleteReq := &pb.DeleteObjectRequest{
		Name: getReq.GetName(),
	}
	deleteRes, err := c.DeleteObject(ctx, deleteReq)
	if err != nil {
		t.Log(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, deleteRes)
}
