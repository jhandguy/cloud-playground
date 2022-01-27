package object

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

type objectServiceServer struct {
	objectPb.ObjectServiceServer
}

func (s *objectServiceServer) CreateObject(_ context.Context, request *objectPb.CreateObjectRequest) (*objectPb.CreateObjectResponse, error) {
	return &objectPb.CreateObjectResponse{Object: request.GetObject()}, nil
}

func (s *objectServiceServer) GetObject(_ context.Context, request *objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
	return &objectPb.GetObjectResponse{Object: &objectPb.Object{
		Id: request.GetId(),
	}}, nil
}

func (s *objectServiceServer) DeleteObject(context.Context, *objectPb.DeleteObjectRequest) (*objectPb.DeleteObjectResponse, error) {
	return &objectPb.DeleteObjectResponse{}, nil
}

func newContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

func newClientConn(t *testing.T) func(ctx context.Context) (*grpc.ClientConn, error) {
	return func(ctx context.Context) (*grpc.ClientConn, error) {
		bufSize := 1024 * 1024
		listener := bufconn.Listen(bufSize)
		bufDialer := func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}

		s := grpc.NewServer()
		objectPb.RegisterObjectServiceServer(s, &objectServiceServer{})
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())

		go func() {
			if err := s.Serve(listener); err != nil {
				t.Error(err)
			}
		}()
		return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}

func TestCreateObject(t *testing.T) {
	api := API{
		CreateObject: CreateObject(newContext, newClientConn(t)),
	}

	req := &objectPb.CreateObjectRequest{
		Object: &objectPb.Object{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}

	ctx := context.Background()
	res, err := api.CreateObject(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.GetObject().GetId(), res.GetObject().GetId())
	assert.Equal(t, req.GetObject().GetContent(), res.GetObject().GetContent())
}

func TestGetObject(t *testing.T) {
	api := API{
		GetObject: GetObject(newContext, newClientConn(t)),
	}

	req := &objectPb.GetObjectRequest{
		Id: uuid.NewString(),
	}

	ctx := context.Background()
	res, err := api.GetObject(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.GetId(), res.GetObject().GetId())
}

func TestDeleteObject(t *testing.T) {
	api := API{
		DeleteObject: DeleteObject(newContext, newClientConn(t)),
	}

	req := &objectPb.DeleteObjectRequest{
		Id: uuid.NewString(),
	}

	ctx := context.Background()
	res, err := api.DeleteObject(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
