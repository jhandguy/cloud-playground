package object

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"

	pb "github.com/jhandguy/devops-playground/s3/pb/object"
)

func TestCreateObject(t *testing.T) {
	var actBucket, actName, actContent string

	api := API{
		S3: S3{
			Bucket: "bucket",
			PutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				actBucket = *input.Bucket
				actName = *input.Key

				byt, err := ioutil.ReadAll(input.Body)
				if err != nil {
					t.Fatal(err)
				}
				actContent = string(byt)

				return &s3.PutObjectOutput{}, nil
			},
			GetObjectWithContext:    nil,
			DeleteObjectWithContext: nil,
		},
	}

	req := &pb.CreateObjectRequest{
		Name:    "name",
		Content: "content",
	}
	resp, err := api.CreateObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.Name, actName)
	assert.Equal(t, req.Content, actContent)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Object.Name)
	assert.Equal(t, req.Content, resp.Object.Content)
}

func TestGetObject(t *testing.T) {
	var actBucket, actName string
	expContent := "content"

	api := API{
		S3: S3{
			Bucket:               "bucket",
			PutObjectWithContext: nil,
			GetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				actBucket = *input.Bucket
				actName = *input.Key

				return &s3.GetObjectOutput{
					Body: ioutil.NopCloser(strings.NewReader(expContent)),
				}, nil
			},
			DeleteObjectWithContext: nil,
		},
	}

	req := &pb.GetObjectRequest{
		Name: "name",
	}
	resp, err := api.GetObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.Name, actName)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Object.Name)
	assert.Equal(t, expContent, resp.Object.Content)
}

func TestDeleteObject(t *testing.T) {
	var actBucket, actName string

	api := API{
		S3: S3{
			Bucket:               "bucket",
			PutObjectWithContext: nil,
			GetObjectWithContext: nil,
			DeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				actBucket = *input.Bucket
				actName = *input.Key

				return &s3.DeleteObjectOutput{}, nil
			},
		},
	}

	req := &pb.DeleteObjectRequest{
		Name: "name",
	}
	resp, err := api.DeleteObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.Name, actName)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
