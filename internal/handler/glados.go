package handler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
)

type Glados struct {
	BaseLogic
	website cfg.Website
}

func init() {
	RegisterCheckInHandler("glados", &Glados{}) // 注册处理器
}

func (i *Glados) getUserInfo() error {
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "GET",
		URL:                "https://glados.network/api/user/status",
		Headers:            i.website.Headers,
		InsecureSkipVerify: true,
		Proxy:              true,
	})
	if err != nil {
		return err
	}
	if code, ok := response["code"].(float64); ok && 0 == code {
		if data, ok := response["data"].(map[string]any); ok {
			i.PushContent("👶 账号：%s", data["email"].(string))
		}
	}
	return nil
}

func (i *Glados) doSign() error {
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                "https://glados.network/api/user/checkin",
		Headers:            i.website.Headers,
		BodyData:           i.website.Body,
		BodyToJson:         true,
		InsecureSkipVerify: true,
		Proxy:              true,
	})
	if err != nil {
		return err
	}
	if code, ok := response["code"].(float64); ok {
		if 0 == code {
			if msg, ok := response["message"].(string); ok {
				i.PushContent("💾 %s", msg)
			}
			if list, ok := response["list"].([]any); ok {
				if list != nil {
					for idx, v := range list {
						lp := v.(map[string]any)
						if idx == 0 {
							balance, err := util.StringToInt(lp["balance"].(string))
							if err != nil {
								return err
							}
							i.PushContent("🎁 当前Points: %d", balance)
						}
					}
				}
			}
		} else {
			if msg, ok := response["message"].(string); ok {
				i.PushContent("🔔 %s", msg)
			}
			if list, ok := response["list"].([]any); ok {
				if list != nil {
					for idx, v := range list {
						lp := v.(map[string]any)
						if idx == 0 {
							balance, err := util.StringToInt(lp["balance"].(string))
							if err != nil {
								return err
							}
							i.PushContent("🎁 当前Points: %d", balance)
						}
					}
				}
			}
		}
	}
	return nil
}

// NewGlados 初始化 Quark 实例
func NewGlados(website cfg.Website) *Glados {
	obj := &Glados{
		website: website,
	}
	obj.Content = "👙 [服务]" + website.Name + "签到信息\n"
	return obj
}

// Run 执行签到操作
func (i *Glados) Run(website cfg.Website) string {
	logger.Log().Debug("----------Glados开始签到----------")
	glados := NewGlados(website)
	_ = glados.getUserInfo()
	err := glados.doSign()
	if err != nil {
		logger.Log().Error("[Glados]签到失败: " + err.Error())
		return ""
	}
	logger.Log().Debug("----------Glados结束签到----------")
	return glados.Content
}
