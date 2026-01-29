package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// newNacosClient creates a Nacos client
func newNacosClient(o options) (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(o.endpoint, o.port),
	}

	cc := constant.ClientConfig{
		NamespaceId:         o.namespaceID,
		TimeoutMs:           o.timeoutMs,
		LogDir:              o.logDir,
		CacheDir:            o.cacheDir,
		LogLevel:            o.logLevel,
		NotLoadCacheAtStart: true,
		Username:            o.username,
		Password:            o.password,
		SecretKey:           o.SecretKey,
		AccessKey:           o.AccessKey,
	}
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
