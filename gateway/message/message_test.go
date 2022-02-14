package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/object"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

func TestCreateMessage(t *testing.T) {
	var actItemReq *itemPb.CreateItemRequest
	var actObjectReq *objectPb.CreateObjectRequest

	expMsg := Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	expItemReq := &itemPb.CreateItemRequest{
		Item: &itemPb.Item{
			Id:      expMsg.ID,
			Content: expMsg.Content,
		},
	}
	expObjectReq := &objectPb.CreateObjectRequest{
		Object: &objectPb.Object{
			Id:      expMsg.ID,
			Content: expMsg.Content,
		},
	}

	api := API{
		ItemAPI: &item.API{
			CreateItem: func(ctx context.Context, req *itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
				actItemReq = req

				return &itemPb.CreateItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
		},
		ObjectAPI: &object.API{
			CreateObject: func(ctx context.Context, req *objectPb.CreateObjectRequest) (*objectPb.CreateObjectResponse, error) {
				actObjectReq = req

				return &objectPb.CreateObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
		},
	}

	byt, err := json.Marshal(expMsg)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/message", api.CreateMessage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/message", bytes.NewReader(byt))

	router.ServeHTTP(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestGetMessage(t *testing.T) {
	var actItemReq *itemPb.GetItemRequest
	var actObjectReq *objectPb.GetObjectRequest

	expMsg := Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	expItemReq := &itemPb.GetItemRequest{
		Id: expMsg.ID,
	}
	expObjectReq := &objectPb.GetObjectRequest{
		Id: expMsg.ID,
	}

	api := API{
		ItemAPI: &item.API{
			GetItem: func(ctx context.Context, req *itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
				actItemReq = req

				return &itemPb.GetItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
		},
		ObjectAPI: &object.API{
			GetObject: func(ctx context.Context, req *objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
				actObjectReq = req

				return &objectPb.GetObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
		},
	}

	byt, err := json.Marshal(expMsg)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/message/:id", api.GetMessage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/message/%s", expMsg.ID), nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestDeleteMessage(t *testing.T) {
	var actItemReq *itemPb.DeleteItemRequest
	var actObjectReq *objectPb.DeleteObjectRequest

	expMsg := Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	expItemReq := &itemPb.DeleteItemRequest{
		Id: expMsg.ID,
	}
	expObjectReq := &objectPb.DeleteObjectRequest{
		Id: expMsg.ID,
	}

	api := API{
		ItemAPI: &item.API{
			DeleteItem: func(ctx context.Context, req *itemPb.DeleteItemRequest) (*itemPb.DeleteItemResponse, error) {
				actItemReq = req

				return &itemPb.DeleteItemResponse{}, nil
			},
		},
		ObjectAPI: &object.API{
			DeleteObject: func(ctx context.Context, req *objectPb.DeleteObjectRequest) (*objectPb.DeleteObjectResponse, error) {
				actObjectReq = req

				return &objectPb.DeleteObjectResponse{}, nil
			},
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/message/:id", api.DeleteMessage)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/message/%s", expMsg.ID), nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, strings.ReplaceAll(w.Body.String(), "\n", ""))
}
