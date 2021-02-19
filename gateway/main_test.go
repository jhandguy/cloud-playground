package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/devops-playground/gateway/item"
	"github.com/jhandguy/devops-playground/gateway/message"
	"github.com/jhandguy/devops-playground/gateway/object"
	itemPb "github.com/jhandguy/devops-playground/gateway/pb/item"
	objectPb "github.com/jhandguy/devops-playground/gateway/pb/object"
)

func TestServeAPI(t *testing.T) {
	var isCreateItemCalled, isCreateObjectCalled bool
	var isGetItemCalled, isGetObjectCalled bool
	var isDeleteItemCalled, isDeleteObjectCalled bool
	var isMiddlewareCalled bool

	expMessage := message.Message{
		ID:      "id",
		Name:    "name",
		Content: "content",
	}

	api := &message.API{
		ItemAPI: &item.API{
			CreateItem: func(*itemPb.CreateItemRequest) (*itemPb.CreateItemResponse, error) {
				isCreateItemCalled = true

				return &itemPb.CreateItemResponse{
					Item: &itemPb.Item{
						Id:      expMessage.ID,
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
			GetItem: func(*itemPb.GetItemRequest) (*itemPb.GetItemResponse, error) {
				isGetItemCalled = true

				return &itemPb.GetItemResponse{
					Item: &itemPb.Item{
						Id:      expMessage.ID,
						Name:    expMessage.Name,
						Content: expMessage.Content,
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
						Name:    expMessage.Name,
						Content: expMessage.Content,
					},
				}, nil
			},
			GetObject: func(*objectPb.GetObjectRequest) (*objectPb.GetObjectResponse, error) {
				isGetObjectCalled = true

				return &objectPb.GetObjectResponse{
					Object: &objectPb.Object{
						Name:    expMessage.Name,
						Content: expMessage.Content,
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

	byt, err := json.Marshal(expMessage)
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
	r = httptest.NewRequest(http.MethodGet, "/message", bytes.NewReader(byt))

	router.ServeHTTP(w, r)

	assert.True(t, isMiddlewareCalled)
	assert.True(t, isGetItemCalled)
	assert.True(t, isGetObjectCalled)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(byt), strings.ReplaceAll(w.Body.String(), "\n", ""))

	isMiddlewareCalled = false
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodDelete, "/message", bytes.NewReader(byt))

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
	expMessage := message.Message{
		Name:    "name",
		Content: "content",
	}

	byt, err := json.Marshal(expMessage)
	if err != nil {
		t.Fatal(err)
	}

	code, body := doRequest(t, url, http.MethodPost, bytes.NewReader(byt))

	assert.Equal(t, http.StatusOK, code)

	var actMessage message.Message
	byt = body.Bytes()
	err = json.Unmarshal(byt, &actMessage)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expMessage.Name, actMessage.Name)
	assert.Equal(t, expMessage.Content, actMessage.Content)

	code, body = doRequest(t, url, http.MethodGet, bytes.NewReader(byt))

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, string(byt), body.String())

	code, body = doRequest(t, url, http.MethodDelete, bytes.NewReader(byt))

	assert.Equal(t, http.StatusOK, code)
	assert.Empty(t, body.String())
}

func doRequest(t *testing.T, url, method string, body io.Reader) (int, *bytes.Buffer) {
	apiKey := retrieveEnv("GATEWAY_API_KEY")
	target := fmt.Sprintf("http://%s/message", url)

	req, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	req.Header.Add("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, buf
}
