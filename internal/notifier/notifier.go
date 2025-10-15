package notifier

import (
	"auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"encoding/json"
	"fmt"
	"net/url"
)

type Notifier struct {
}

func New() *Notifier {
	return &Notifier{}
}

func (n *Notifier) SendWeCom(message string) error {
	if config.Cfg.Notifications.WeCom.Webhook == "" {
		logger.Log().Debug("未配置企微消息推送")
		return nil
	}
	logger.Log().Debug("开始执行企微消息推送")
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                config.Cfg.Notifications.WeCom.Webhook,
		QueryParams:        nil,
		BodyData:           jsonPayload,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	errcode, ok := resp["errcode"].(float64)
	if !ok {
		logger.Log().Error("企微消息推送失败: 无效的errcode类型")
		return fmt.Errorf("企微消息推送失败: 无效的errcode类型")
	}
	if errcode != 0 {
		errmsg, ok := resp["errmsg"].(string)
		if !ok {
			logger.Log().Error("企微消息推送失败: 无效的errmsg类型")
			return fmt.Errorf("企微消息推送失败: 无效的errmsg类型")
		}
		logger.Log().Errorf("企微消息推送失败: %s", errmsg)
		return fmt.Errorf("企微消息推送失败: %s", errmsg)
	}
	return nil
}

func (n *Notifier) SendTelegram(message string) error {
	if config.Cfg.Notifications.Telegram.BotToken == "" || config.Cfg.Notifications.Telegram.UID == "" {
		logger.Log().Debug("tg 服务的 bot_token 或者 user_id 未设置!!")
		return nil
	}
	logger.Log().Info("开始执行Telegram消息推送")
	var apiUrl string
	if config.Cfg.Notifications.Telegram.APIHost != "" {
		apiUrl = fmt.Sprintf("https://%s/bot%s/sendMessage", config.Cfg.Notifications.Telegram.APIHost, config.Cfg.Notifications.Telegram.BotToken)
	} else {
		apiUrl = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Cfg.Notifications.Telegram.BotToken)
	}
	logger.Log().Debug("开始拼装telegram消息推送参数")
	formData := url.Values{}
	formData.Add("chat_id", config.Cfg.Notifications.Telegram.UID)
	formData.Add("text", message)
	formData.Add("disable_web_page_preview", "true")

	var proxyURL string
	if config.Cfg.Proxy.Host != "" && config.Cfg.Proxy.Port != "" {
		proxyURL = fmt.Sprintf("%s:%s", config.Cfg.Proxy.Host, config.Cfg.Proxy.Port)
	}
	logger.Log().Info("开始执行Telegram消息推送")
	result, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                apiUrl,
		BodyData:           formData,
		InsecureSkipVerify: false,
		Proxy:              proxyURL,
	})
	if err != nil {
		logger.Log().Errorf("telegram消息推送失败: %v", err)
		return fmt.Errorf("telegram消息推送失败: %v", err)
	}
	if ok, exists := result["ok"].(bool); exists && ok {
		logger.Log().Info("telegram推送成功！")
	} else {
		logger.Log().Error("telegram推送失败！")
	}
	return nil
}
