package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("%v != %v", a, b))
}

func RecordRequest(fun func(http.ResponseWriter, *http.Request), method, name string, body io.Reader) (int, *bytes.Buffer) {
	target := "/object"
	if len(name) > 0 {
		target = fmt.Sprintf("%s?name=%s", target, name)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, target, body)

	fun(w, r)

	return w.Code, w.Body
}

func SendRequest(t *testing.T, url, method, name, apiKey string, body io.Reader) (int, *bytes.Buffer) {
	target := fmt.Sprintf("%s/object", url)
	if len(name) > 0 {
		target = fmt.Sprintf("%s?name=%s", target, name)
	}

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

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, buf
}
