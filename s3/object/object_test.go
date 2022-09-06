package object

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pb "github.com/jhandguy/cloud-playground/s3/pb/object"
)

func TestCreateObject(t *testing.T) {
	var actBucket, actID, actContent string

	api := API{
		S3: S3{
			Bucket: "bucket",
			PutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				actBucket = *input.Bucket
				actID = *input.Key

				byt, err := io.ReadAll(input.Body)
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
		Object: &pb.Object{
			Id:      uuid.NewString(),
			Content: "content",
		},
	}
	resp, err := api.CreateObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.GetObject().GetId(), actID)
	assert.Equal(t, req.GetObject().GetContent(), actContent)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.GetObject().GetId(), resp.Object.GetId())
	assert.Equal(t, req.GetObject().GetContent(), resp.Object.GetContent())
}

func TestGetObject(t *testing.T) {
	var actBucket, actID string
	expContent := "content"

	api := API{
		S3: S3{
			Bucket:               "bucket",
			PutObjectWithContext: nil,
			GetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				actBucket = *input.Bucket
				actID = *input.Key

				return &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader(expContent)),
				}, nil
			},
			DeleteObjectWithContext: nil,
		},
	}

	req := &pb.GetObjectRequest{
		Id: uuid.NewString(),
	}
	resp, err := api.GetObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.GetId(), actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.GetId(), resp.Object.GetId())
	assert.Equal(t, expContent, resp.Object.GetContent())
}

func TestDeleteObject(t *testing.T) {
	var actBucket, actID string

	api := API{
		S3: S3{
			Bucket:               "bucket",
			PutObjectWithContext: nil,
			GetObjectWithContext: nil,
			DeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				actBucket = *input.Bucket
				actID = *input.Key

				return &s3.DeleteObjectOutput{}, nil
			},
		},
	}

	req := &pb.DeleteObjectRequest{
		Id: uuid.NewString(),
	}
	resp, err := api.DeleteObject(context.Background(), req)

	assert.Equal(t, api.S3.Bucket, actBucket)
	assert.Equal(t, req.GetId(), actID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
