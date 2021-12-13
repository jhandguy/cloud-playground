package item

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
)

type itemServiceServer struct {
	itemPb.ItemServiceServer
}

func (s *itemServiceServer) CreateItem(_ context.Context, request *itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
	return &itemPb.CreateItemResponse{Item: request.GetItem()}, nil
}

func (s *itemServiceServer) GetItem(_ context.Context, request *itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
	return &itemPb.GetItemResponse{Item: &itemPb.Item{
		Id: request.GetId(),
	}}, nil
}

func (s *itemServiceServer) DeleteItem(context.Context, *itemPb.DeleteItemRequest) (*itemPb.DeleteItemResponse, error) {
	return &itemPb.DeleteItemResponse{}, nil
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
		itemPb.RegisterItemServiceServer(s, &itemServiceServer{})
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())

		go func() {
			if err := s.Serve(listener); err != nil {
				t.Error(err)
			}
		}()
		return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	}
}

func TestCreateItem(t *testing.T) {
	api := API{
		CreateItem: CreateItem(newContext, newClientConn(t)),
	}

	req := &itemPb.CreateItemRequest{
		Item: &itemPb.Item{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}

	ctx := context.Background()
	res, err := api.CreateItem(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.GetItem().GetId(), res.GetItem().GetId())
	assert.Equal(t, req.GetItem().GetContent(), res.GetItem().GetContent())
}

func TestGetItem(t *testing.T) {
	api := API{
		GetItem: GetItem(newContext, newClientConn(t)),
	}

	req := &itemPb.GetItemRequest{
		Id: uuid.NewString(),
	}

	ctx := context.Background()
	res, err := api.GetItem(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.GetId(), res.GetItem().GetId())
}

func TestDeleteItem(t *testing.T) {
	api := API{
		DeleteItem: DeleteItem(newContext, newClientConn(t)),
	}

	req := &itemPb.DeleteItemRequest{
		Id: uuid.NewString(),
	}

	ctx := context.Background()
	res, err := api.DeleteItem(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
