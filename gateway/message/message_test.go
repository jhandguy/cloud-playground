package message

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/object"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

func TestCreateMessage(t *testing.T) {
	var actItemReq *itemPb.CreateItemRequest
	var actObjectReq *objectPb.CreateObjectRequest

	expMessage := Message{
		ID:      "id",
		Name:    "name",
		Content: "content",
	}

	expItemReq := &itemPb.CreateItemRequest{
		Name:    expMessage.Name,
		Content: expMessage.Content,
	}
	expObjectReq := &objectPb.CreateObjectRequest{
		Name:    expMessage.Name,
		Content: expMessage.Content,
	}

	api := API{
		ItemAPI: &item.API{
			CreateItem: func(req *itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
				actItemReq = req

				return &itemPb.CreateItemResponse{
					Item: &itemPb.Item{
						Id:      expMessage.ID,
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
		},
		ObjectAPI: &object.API{
			CreateObject: func(req *objectPb.CreateObjectRequest) (*objectPb.CreateObjectResponse, error) {
				actObjectReq = req

				return &objectPb.CreateObjectResponse{
					Object: &objectPb.Object{
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
		},
	}

	byt, err := json.Marshal(expMessage)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/message", bytes.NewReader(byt))

	api.CreateMessage(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestGetMessage(t *testing.T) {
	var actItemReq *itemPb.GetItemRequest
	var actObjectReq *objectPb.GetObjectRequest

	expMessage := Message{
		ID:      "id",
		Name:    "name",
		Content: "content",
	}

	expItemReq := &itemPb.GetItemRequest{
		Id: expMessage.ID,
	}
	expObjectReq := &objectPb.GetObjectRequest{
		Name: expMessage.Name,
	}

	api := API{
		ItemAPI: &item.API{
			GetItem: func(req *itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
				actItemReq = req

				return &itemPb.GetItemResponse{
					Item: &itemPb.Item{
						Id:      expMessage.ID,
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
		},
		ObjectAPI: &object.API{
			GetObject: func(req *objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
				actObjectReq = req

				return &objectPb.GetObjectResponse{
					Object: &objectPb.Object{
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
		},
	}

	byt, err := json.Marshal(expMessage)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/message", bytes.NewReader(byt))

	api.GetMessage(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestDeleteMessage(t *testing.T) {
	var actItemReq *itemPb.DeleteItemRequest
	var actObjectReq *objectPb.DeleteObjectRequest

	expMessage := Message{
		ID:      "id",
		Name:    "name",
		Content: "content",
	}

	expItemReq := &itemPb.DeleteItemRequest{
		Id: expMessage.ID,
	}
	expObjectReq := &objectPb.DeleteObjectRequest{
		Name: expMessage.Name,
	}

	api := API{
		ItemAPI: &item.API{
			DeleteItem: func(req *itemPb.DeleteItemRequest) (*itemPb.DeleteItemResponse, error) {
				actItemReq = req

				return &itemPb.DeleteItemResponse{}, nil
			},
		},
		ObjectAPI: &object.API{
			DeleteObject: func(req *objectPb.DeleteObjectRequest) (*objectPb.DeleteObjectResponse, error) {
				actObjectReq = req

				return &objectPb.DeleteObjectResponse{}, nil
			},
		},
	}

	byt, err := json.Marshal(expMessage)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/message", bytes.NewReader(byt))

	api.DeleteMessage(w, r)

	assert.Equal(t, expItemReq, actItemReq)
	assert.Equal(t, expObjectReq, actObjectReq)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, strings.ReplaceAll(w.Body.String(), "\n", ""))
}
