package driver

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type COS struct {
	Client *cos.Client
}

func NewCOS(bucketURL, region, secretID, secretKey string) *COS {
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &COS{Client: client}
}

func (c *COS) Put(ctx context.Context, key string, r io.Reader) (string, error) {
	_, err := c.Client.Object.Put(ctx, key, r, nil)
	if err != nil {
		return "", err
	}

	return c.Client.BaseURL.BucketURL.String() + "/" + key, nil
}

func (c *COS) Delete(ctx context.Context, key string) error {
	_, err := c.Client.Object.Delete(ctx, key)
	return err
}
