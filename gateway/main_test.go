package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

func TestIsValidAPIKey(t *testing.T) {
	apiKey := "api-key"

	auth := fmt.Sprintf("Bearer %s", apiKey)
	assert.True(t, isValidAPIKey(auth, apiKey))

	auth = ""
	assert.False(t, isValidAPIKey(auth, apiKey))

	auth = apiKey
	assert.True(t, isValidAPIKey(auth, apiKey))

	auth = "wrong"
	assert.False(t, isValidAPIKey(auth, apiKey))
}

func TestServeAPI(t *testing.T) {
	var isCreateItemCalled, isCreateObjectCalled bool
	var isGetItemCalled, isGetObjectCalled bool
	var isDeleteItemCalled, isDeleteObjectCalled bool
	var isMiddlewareCalled bool

	expMsg := message.Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	api := &message.API{
		ItemAPI: &item.API{
			CreateItem: func(*itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
				isCreateItemCalled = true

				return &itemPb.CreateItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			GetItem: func(*itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
				isGetItemCalled = true

				return &itemPb.GetItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			DeleteItem: func(*itemPb.DeleteItemRequest) (*itemPb.DeleteItemResponse, error) {
				isDeleteItemCalled = true

				return &itemPb.DeleteItemResponse{}, nil
			},
		},
		ObjectAPI: &object.API{
			CreateObject: func(*objectPb.CreateObjectRequest) (*objectPb.CreateObjectResponse, error) {
				isCreateObjectCalled = true

				return &objectPb.CreateObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			GetObject: func(*objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
				isGetObjectCalled = true

				return &objectPb.GetObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			DeleteObject: func(*objectPb.DeleteObjectRequest) (*objectPb.DeleteObjectResponse, error) {
				isDeleteObjectCalled = true

				return &objectPb.DeleteObjectResponse{}, nil
			},
		},
	}

	middleware := func(next http.Handler) http.Handler {
		isMiddlewareCalled = true
		return next
	}

	router := serveAPI(api, middleware)

	byt, err := json.Marshal(expMsg)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/message", bytes.NewReader(byt))

	router.ServeHTTP(w, r)

	assert.True(t, isMiddlewareCalled)
	assert.True(t, isCreateItemCalled)
	assert.True(t, isCreateObjectCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))

	isMiddlewareCalled = false
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/message/%s", expMsg.ID), nil)

	router.ServeHTTP(w, r)

	assert.True(t, isMiddlewareCalled)
	assert.True(t, isGetItemCalled)
	assert.True(t, isGetObjectCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))

	isMiddlewareCalled = false
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/message/%s", expMsg.ID), nil)

	router.ServeHTTP(w, r)

	assert.True(t, isMiddlewareCalled)
	assert.True(t, isDeleteItemCalled)
	assert.True(t, isDeleteObjectCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, strings.ReplaceAll(w.Body.String(), "\n", ""))
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	go main()

	port := retrieveEnv("GATEWAY_PORT")
	url := fmt.Sprintf("localhost:%s", port)
	testGateway(t, url)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := retrieveEnv("GATEWAY_URL")
	testGateway(t, url)
}

func testGateway(t *testing.T, url string) {
	apiKey := retrieveEnv("GATEWAY_API_KEY")
	client := resty.
		New().
		SetHostURL(fmt.Sprintf("http://%s", url)).
		SetAuthToken(apiKey).
		SetDebug(true)

	expMsg := message.Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	res, err := client.
		R().
		SetResult(message.Message{}).
		SetBody(expMsg).
		Post("/message")
	if err != nil {
		t.Fatal(err)
	}

	actMsg := res.Result().(*message.Message)

	assert.Equal(t, http.StatusOK, res.StatusCode())
	assert.Equal(t, expMsg.ID, actMsg.ID)
	assert.Equal(t, expMsg.Content, actMsg.Content)

	res, err = client.
		R().
		SetResult(message.Message{}).
		SetPathParam("id", expMsg.ID).
		Get("/message/{id}")
	if err != nil {
		t.Fatal(err)
	}

	actMsg = res.Result().(*message.Message)

	assert.Equal(t, http.StatusOK, res.StatusCode())
	assert.Equal(t, expMsg.ID, actMsg.ID)
	assert.Equal(t, expMsg.Content, actMsg.Content)

	res, err = client.
		R().
		SetPathParam("id", expMsg.ID).
		Delete("/message/{id}")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode())
	assert.Empty(t, res.Body())
}
