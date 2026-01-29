package doubao

import (
	"sync"
	"time"

	"github.com/iWuxc/go-wit/log"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

var (
	client *arkruntime.Client
	mu     sync.Mutex
)

// Config 豆包配置
type Config struct {
	ApiKey      string `json:"api_key" mapstructure:"api_key"`
	OutfitModel string `json:"outfit_model" mapstructure:"outfit_model"`
	ImageModel  string `json:"image_model" mapstructure:"image_model"`
}

// InitConfig 初始化豆包配置
func InitConfig(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	// 验证配置
	if config.ApiKey == "" {
		return nil // 如果ApiKey为空，不初始化client
	}

	// 创建client
	client = arkruntime.NewClientWithApiKey(
		config.ApiKey,
		arkruntime.WithTimeout(30*time.Minute),
	)

	log.Infof("豆包服务初始化成功，ApiKey=%s, OutfitModel=%s, ImageModel=%s",
		maskString(config.ApiKey),
		config.OutfitModel,
		config.ImageModel,
	)

	return nil
}

// GetClient 获取豆包客户端
func GetClient() *arkruntime.Client {
	mu.Lock()
	defer mu.Unlock()
	return client
}

// ResetClient 重置客户端（例如在配置变更时调用）
func ResetClient() {
	mu.Lock()
	defer mu.Unlock()
	client = nil
}

// maskString 掩盖敏感信息
func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
