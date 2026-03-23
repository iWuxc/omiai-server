package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/iWuxc/go-wit/log"
)

// Config 消息通知配置
type Config struct {
	WecomWebhookURL string // 企业微信机器人 Webhook
	WeappTemplateID string // 小程序订阅消息模板ID
}

// Service 通知服务接口
type Service interface {
	// NotifyManager 通知红娘 (如：新线索、互相心动)
	NotifyManager(ctx context.Context, managerID uint64, title, content string) error
	
	// NotifyClient 通知C端用户 (如：匹配成功、状态变更)
	NotifyClient(ctx context.Context, clientID uint64, openID, title, content string) error
}

type NotificationService struct {
	config *Config
}

func NewNotificationService() Service {
	// TODO: 从配置中心加载真实配置
	return &NotificationService{
		config: &Config{
			WecomWebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=mock_key",
		},
	}
}

func (s *NotificationService) NotifyManager(ctx context.Context, managerID uint64, title, content string) error {
	// 模拟发送企业微信消息
	// 实际生产中应调用企业微信 API，通过 managerID 找到对应的企微 userID，然后发送应用消息
	log.Infof("[企微通知 -> 红娘 %d] 标题: %s, 内容: %s, 时间: %s", 
		managerID, title, content, time.Now().Format("2006-01-02 15:04:05"))
		
	// TODO: 发送 HTTP 请求到企微 Webhook 或应用接口
	// http.Post(s.config.WecomWebhookURL, "application/json", body)
	
	return nil
}

func (s *NotificationService) NotifyClient(ctx context.Context, clientID uint64, openID, title, content string) error {
	if openID == "" {
		return fmt.Errorf("client %d has no openid", clientID)
	}
	
	// 模拟发送微信小程序订阅消息 / 公众号模板消息
	log.Infof("[微信通知 -> 客户 %d (OpenID: %s)] 标题: %s, 内容: %s", 
		clientID, openID, title, content)
		
	// TODO: 调用微信服务端 API 发送订阅消息
	
	return nil
}
