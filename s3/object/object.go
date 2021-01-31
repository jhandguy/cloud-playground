package object

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type API struct {
	S3 S3
	ObjectServiceServer
}

type S3 struct {
	Bucket string

	PutObjectWithContext    func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
	GetObjectWithContext    func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	DeleteObjectWithContext func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error)
}

func (api *API) CreateObject(ctx context.Context, req *CreateObjectRequest) (*CreateObjectResponse, error) {
	_, err := api.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(req.Name),
		Body:   strings.NewReader(req.Content),
	})
	if err != nil {
		return nil, err
	}

	return &CreateObjectResponse{
		Object: &Object{
			Name:    req.Name,
			Content: req.Content,
		},
	}, nil
}

func (api *API) GetObject(ctx context.Context, req *GetObjectRequest) (*GetObjectResponse, error) {
	out, err := api.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(req.Name),
	})
	if err != nil {
		return nil, err
	}

	body := out.Body
	if body == nil {
		return &GetObjectResponse{}, nil
	}

	defer func() {
		_ = body.Close()
	}()

	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return &GetObjectResponse{
		Object: &Object{
			Name:    req.Name,
			Content: string(byt),
		},
	}, nil
}

func (api *API) DeleteObject(ctx context.Context, req *DeleteObjectRequest) (*DeleteObjectResponse, error) {
	_, err := api.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(api.S3.Bucket),
		Key:    aws.String(req.Name),
	})

	return &DeleteObjectResponse{}, err
}
