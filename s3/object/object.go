package object

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Object struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateObject(createObject func(string, Object) error, bucket string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj Object

		err := json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = createObject(bucket, obj)
		if err != nil {
			log.Printf("failed to upload object %s in bucket %s\nwith content: %s\nerror: %v\n", obj.Name, bucket, obj.Content, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(obj)
		if err != nil {
			log.Printf("failed to encode object: %v\nerror: %v\n", obj, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("successfully uploaded object %s in bucket %s\nwith content: %s\n", obj.Name, bucket, obj.Content)
	}
}

func GetObject(getObject func(string, string) (io.ReadCloser, error), bucket string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := getNameQuery(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := getObject(bucket, name)
		if err != nil {
			log.Printf("failed to download object %s in bucket %s\nerror: %v\n", name, bucket, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer func() {
			_ = body.Close()
		}()

		byt, err := ioutil.ReadAll(body)
		if err != nil {
			log.Printf("failed to read object: %v\nerror: %v\n", body, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		obj := Object{
			Name:    name,
			Content: string(byt),
		}

		err = json.NewEncoder(w).Encode(obj)
		if err != nil {
			log.Printf("failed to encode object: %v\nerror: %v\n", obj, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("successfully downloaded object %s in bucket %s\nwith content: %s\n", name, bucket, obj.Content)
	}
}

func DeleteObject(deleteObject func(string, string) error, bucket string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := getNameQuery(r)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := deleteObject(bucket, name)
		if err != nil {
			log.Printf("failed to delete object %s in bucket %s\nerror: %v\n", name, bucket, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("successfully deleted object %s in bucket %s\n", name, bucket)
	}
}

func CreateObjectFunc(client *s3.S3) func(string, Object) error {
	return func(bucket string, obj Object) error {
		_, err := client.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(obj.Name),
			Body:   strings.NewReader(obj.Content),
		})

		return err
	}
}

func GetObjectFunc(client *s3.S3) func(string, string) (io.ReadCloser, error) {
	return func(bucket string, name string) (io.ReadCloser, error) {
		out, err := client.GetObjectWithContext(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(name),
		})
		if err != nil {
			return nil, err
		}

		return out.Body, nil
	}
}

func DeleteObjectFunc(client *s3.S3) func(string, string) error {
	return func(bucket string, name string) error {
		_, err := client.DeleteObjectWithContext(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(name),
		})

		return err
	}
}

func getNameQuery(r *http.Request) (string, bool) {
	query, ok := r.URL.Query()["name"]

	if !ok || len(query[0]) < 1 {
		return "", false
	}

	return query[0], true
}
