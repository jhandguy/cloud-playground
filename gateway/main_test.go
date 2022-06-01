package main

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
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/cloud-playground/gateway/item"
	"github.com/jhandguy/cloud-playground/gateway/message"
	"github.com/jhandguy/cloud-playground/gateway/object"
	itemPb "github.com/jhandguy/cloud-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/cloud-playground/gateway/pb/object"
)

func TestIsValidToken(t *testing.T) {
	token := "token"

	auth := fmt.Sprintf("Bearer %s", token)
	assert.True(t, isValidToken(auth, token))

	auth = ""
	assert.False(t, isValidToken(auth, token))

	auth = token
	assert.True(t, isValidToken(auth, token))

	auth = "wrong"
	assert.False(t, isValidToken(auth, token))
}

func TestRouteAPI(t *testing.T) {
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
			CreateItem: func(context.Context, *itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
				isCreateItemCalled = true

				return &itemPb.CreateItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			GetItem: func(context.Context, *itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
				isGetItemCalled = true

				return &itemPb.GetItemResponse{
					Item: &itemPb.Item{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			DeleteItem: func(context.Context, *itemPb.DeleteItemRequest) (*itemPb.DeleteItemResponse, error) {
				isDeleteItemCalled = true

				return &itemPb.DeleteItemResponse{}, nil
			},
		},
		ObjectAPI: &object.API{
			CreateObject: func(context.Context, *objectPb.CreateObjectRequest) (*objectPb.CreateObjectResponse, error) {
				isCreateObjectCalled = true

				return &objectPb.CreateObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			GetObject: func(context.Context, *objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
				isGetObjectCalled = true

				return &objectPb.GetObjectResponse{
					Object: &objectPb.Object{
						Id:      expMsg.ID,
						Content: expMsg.Content,
					},
				}, nil
			},
			DeleteObject: func(context.Context, *objectPb.DeleteObjectRequest) (*objectPb.DeleteObjectResponse, error) {
				isDeleteObjectCalled = true

				return &objectPb.DeleteObjectResponse{}, nil
			},
		},
	}

	middleware := func(c *gin.Context) {
		isMiddlewareCalled = true
		c.Next()
	}

	gin.SetMode(gin.TestMode)
	router := routeAPI(api, middleware)

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

	port := viper.GetString("gateway-http-port")
	url := fmt.Sprintf("localhost:%s", port)
	testGateway(t, url)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := viper.GetString("gateway-url")
	testGateway(t, url)
}

func testGateway(t *testing.T, url string) {
	token := viper.GetString("gateway-token")
	host := viper.GetString("gateway-host")
	client := resty.
		New().
		SetBaseURL(fmt.Sprintf("http://%s", url)).
		SetAuthToken(token).
		SetHeader("Host", host).
		SetDebug(true).
		SetRetryCount(3)

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
