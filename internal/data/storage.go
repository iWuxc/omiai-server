package data

import (
	"omiai-server/internal/conf"
	"omiai-server/pkg/storage"
	"omiai-server/pkg/storage/driver"
)

func NewStorage(c *conf.Config) (storage.Driver, error) {
	s := c.Storage
	if s == nil || s.Driver == "" || s.Driver == "local" {
		return driver.NewLocal(c.Runtime.Path, ""), nil
	}

	switch s.Driver {
	case "oss":
		return driver.NewOSS(s.OSS.Endpoint, s.OSS.AccessKeyID, s.OSS.AccessKeySecret, s.OSS.BucketName, s.OSS.Domain)
	case "cos":
		return driver.NewCOS(s.COS.BucketURL, s.COS.Region, s.COS.SecretID, s.COS.SecretKey), nil
	default:
		return driver.NewLocal(c.Runtime.Path, ""), nil
	}
}
