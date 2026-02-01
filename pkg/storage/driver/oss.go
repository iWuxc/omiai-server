package driver

import (
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	Client *oss.Client
	Bucket *oss.Bucket
	Domain string
}

func NewOSS(endpoint, accessKeyID, accessKeySecret, bucketName, domain string) (*OSS, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &OSS{
		Client: client,
		Bucket: bucket,
		Domain: domain,
	}, nil
}

func (o *OSS) Put(ctx context.Context, key string, r io.Reader) (string, error) {
	err := o.Bucket.PutObject(key, r)
	if err != nil {
		return "", err
	}

	return o.Domain + "/" + key, nil
}

func (o *OSS) Delete(ctx context.Context, key string) error {
	return o.Bucket.DeleteObject(key)
}
