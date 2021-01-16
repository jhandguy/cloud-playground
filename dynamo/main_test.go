package main

import (
	"bytes"
	"dynamo/item"
	"dynamo/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleItem(t *testing.T) {
	isCreateItemCalled := false
	isGetItemCalled := false
	isDeleteItemCalled := false

	createItem := func(http.ResponseWriter, *http.Request) {
		isCreateItemCalled = true
	}
	getItem := func(http.ResponseWriter, *http.Request) {
		isGetItemCalled = true
	}
	deleteItem := func(http.ResponseWriter, *http.Request) {
		isDeleteItemCalled = true
	}

	handleItem := handleItem(createItem, getItem, deleteItem)

	_, _ = test.RecordRequest(handleItem, http.MethodPost, "", nil)
	test.AssertEqual(t, isCreateItemCalled, true)

	_, _ = test.RecordRequest(handleItem, http.MethodGet, "", nil)
	test.AssertEqual(t, isGetItemCalled, true)

	_, _ = test.RecordRequest(handleItem, http.MethodDelete, "", nil)
	test.AssertEqual(t, isDeleteItemCalled, true)
}

func TestAuthMiddleware(t *testing.T) {
	isNextFuncCalled := false

	nextFunc := func(http.ResponseWriter, *http.Request) {
		isNextFuncCalled = true
	}

	expApiKey := "1234"

	authMiddlewareFunc := authMiddleware(expApiKey, nextFunc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("Authorization", "wrong")

	authMiddlewareFunc(w, r)

	test.AssertEqual(t, w.Code, http.StatusForbidden)
	test.AssertEqual(t, isNextFuncCalled, false)

	r.Header.Set("Authorization", expApiKey)

	authMiddlewareFunc(w, r)

	test.AssertEqual(t, isNextFuncCalled, true)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	itm := item.Item{
		Id:      "id",
		Name:    "item",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(itm)

	handleItem := handleItemFunc()

	code, _ := test.RecordRequest(handleItem, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, body := test.RecordRequest(handleItem, http.MethodPost, "", bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, body = test.RecordRequest(handleItem, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.RecordRequest(handleItem, http.MethodDelete, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.RecordRequest(handleItem, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := retrieveEnv("DYNAMO_URL")
	apiKey := retrieveEnv("DYNAMO_API_KEY")

	itm := item.Item{
		Id:      "id",
		Name:    "item",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(itm)

	code, _ := test.SendRequest(t, url, http.MethodGet, itm.Id, apiKey, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, body := test.SendRequest(t, url, http.MethodPost, "", apiKey, bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, body = test.SendRequest(t, url, http.MethodGet, itm.Id, apiKey, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.SendRequest(t, url, http.MethodDelete, itm.Id, apiKey, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.SendRequest(t, url, http.MethodGet, itm.Id, apiKey, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}
