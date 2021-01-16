package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	obj := object.Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	handleObject := handleObjectFunc()

	code, _ := test.RecordRequest(handleObject, http.MethodGet, obj.Name, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, body := test.RecordRequest(handleObject, http.MethodPost, "", bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, body = test.RecordRequest(handleObject, http.MethodGet, obj.Name, nil)
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
	apiKey := retrieveEnv("S3_API_KEY")

	obj := object.Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	code, _ := test.SendRequest(t, url, http.MethodGet, obj.Name, apiKey, nil)
	test.AssertEqual(t, code, http.StatusNotFound)

	code, body := test.SendRequest(t, url, http.MethodPost, "", apiKey, bytes.NewReader(byt))
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, body = test.SendRequest(t, url, http.MethodGet, obj.Name, apiKey, nil)
	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))

	code, _ = test.SendRequest(t, url, http.MethodDelete, obj.Name, apiKey, nil)
	test.AssertEqual(t, code, http.StatusOK)

	code, _ = test.SendRequest(t, url, http.MethodGet, obj.Name, apiKey, nil)
	test.AssertEqual(t, code, http.StatusNotFound)
}
