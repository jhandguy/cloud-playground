package main

import (
	"bytes"
	"dynamo/item"
	"dynamo/test"
	"encoding/json"
	"net/http"
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

	itm := item.Item{
		Id:      "id",
		Name:    "item",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(itm)

	code, _ := test.SendRequest(t, url, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, body := test.SendRequest(t, url, http.MethodPost, "", bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, body = test.SendRequest(t, url, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.SendRequest(t, url, http.MethodDelete, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.SendRequest(t, url, http.MethodGet, itm.Id, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}
