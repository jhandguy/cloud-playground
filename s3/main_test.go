package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"s3/object"
	"s3/test"
	"strings"
	"testing"
)

func TestHandleObject(t *testing.T) {
	isCreateObjectCalled := false
	isGetObjectCalled := false
	isDeleteObjectCalled := false

	createObject := func(http.ResponseWriter, *http.Request) {
		isCreateObjectCalled = true
	}
	getObject := func(http.ResponseWriter, *http.Request) {
		isGetObjectCalled = true
	}
	deleteObject := func(http.ResponseWriter, *http.Request) {
		isDeleteObjectCalled = true
	}

	handleObject := handleObject(createObject, getObject, deleteObject)

	_, _ = test.RecordRequest(handleObject, http.MethodPost, "", nil)
	test.AssertEqual(t, isCreateObjectCalled, true)

	_, _ = test.RecordRequest(handleObject, http.MethodGet, "", nil)
	test.AssertEqual(t, isGetObjectCalled, true)

	_, _ = test.RecordRequest(handleObject, http.MethodDelete, "", nil)
	test.AssertEqual(t, isDeleteObjectCalled, true)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	obj := object.Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	handleObject := handleObjectFunc()

	code, _ := test.RecordRequest(handleObject, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, _ = test.RecordRequest(handleObject, http.MethodPost, "", bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)

	code, body := test.RecordRequest(handleObject, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.RecordRequest(handleObject, http.MethodDelete, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.RecordRequest(handleObject, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}

func TestSystem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	url := retrieveEnv("S3_URL")

	obj := object.Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	code, _ := test.SendRequest(t, url, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, _ = test.SendRequest(t, url, http.MethodPost, "", bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)

	code, body := test.SendRequest(t, url, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.SendRequest(t, url, http.MethodDelete, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.SendRequest(t, url, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}
