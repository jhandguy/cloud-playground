package object

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"s3/test"
	"strings"
	"testing"
)

func TestCreateObject(t *testing.T) {
	var actBucket string
	var actObj Object

	fun := func(bucket string, obj Object) error {
		actBucket = bucket
		actObj = obj
		return nil
	}

	bucket := "bucket"
	sut := CreateObject(fun, bucket)

	obj := Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	code, _ := test.RecordRequest(sut, http.MethodPost, "", bytes.NewReader(byt))

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actBucket, bucket)
	test.AssertEqual(t, actObj, obj)
}

func TestGetObject(t *testing.T) {
	var actBucket string
	var actName string

	obj := Object{
		Name:    "obj",
		Content: "Hello world!",
	}
	byt, _ := json.Marshal(obj)

	fun := func(bucket string, name string) (io.ReadCloser, error) {
		actBucket = bucket
		actName = name

		return ioutil.NopCloser(strings.NewReader(obj.Content)), nil
	}

	bucket := "bucket"
	name := "obj"
	sut := GetObject(fun, bucket)

	code, body := test.RecordRequest(sut, http.MethodGet, name, nil)

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actBucket, bucket)
	test.AssertEqual(t, actName, name)
	test.AssertEqual(t, strings.ReplaceAll(body.String(), "\n", ""), string(byt))
}

func TestDeleteObject(t *testing.T) {
	var actBucket string
	var actName string

	fun := func(bucket string, name string) error {
		actBucket = bucket
		actName = name

		return nil
	}

	bucket := "bucket"
	name := "obj"
	sut := DeleteObject(fun, bucket)

	code, _ := test.RecordRequest(sut, http.MethodDelete, name, nil)

	test.AssertEqual(t, code, http.StatusOK)
	test.AssertEqual(t, actBucket, bucket)
	test.AssertEqual(t, actName, name)
}
