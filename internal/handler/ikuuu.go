package handler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
)

type Ikuuu struct {
	BaseLogic
	Headers map[string]string
}

func init() {
	RegisterCheckInHandler("ikuuu", &Ikuuu{}) // 注册处理器
}

func (i *Ikuuu) doSign() error {
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                "https://ikuuu.de/user/checkin",
		Headers:            i.Headers,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	i.PushContent("💾 %s", response["msg"].(string))
	return nil
}

// NewIkuuu 初始化 Quark 实例
func NewIkuuu(website cfg.Website) *Ikuuu {
	obj := &Ikuuu{
		Headers: website.Headers,
	}
	obj.Content = "👙 [服务]" + website.Name + "签到信息\n"
	return obj
}

// Run 执行签到操作
func (i *Ikuuu) Run(website cfg.Website) string {
	logger.Log().Debug("----------IKuuu开始签到----------")
	ikuuu := NewIkuuu(website)
	err := ikuuu.doSign()
	if err != nil {
		logger.Log().Error("[ikuuu]签到失败: " + err.Error())
		return ""
	}
	logger.Log().Debug("----------IKuuu结束签到----------")
	return ikuuu.Content
}
