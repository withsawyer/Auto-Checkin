package scheduler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
)

type Ikuuu struct {
	Headers map[string]string
}

func init() {
	RegisterCheckInHandler("quark", &Quark{}) // 注册处理器
}

func (i *Ikuuu) doSign() (string, error) {
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                "https://ikuuu.de/user/checkin",
		Headers:            i.Headers,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	return response["msg"].(string), nil
}

// NewIkuuu 初始化 Quark 实例
func NewIkuuu(website cfg.Website) *Ikuuu {
	return &Ikuuu{
		Headers: website.Headers,
	}
}

// Run 执行签到操作
func (i *Ikuuu) Run(website cfg.Website) string {
	logger.Log().Info("----------IKuuu开始签到----------")
	ikuuu := NewIkuuu(website)
	msg, err := ikuuu.doSign()
	if err != nil {
		logger.Log().Error("[ikuuu]签到失败: " + err.Error())
		return "签到失败: " + err.Error()
	}
	logger.Log().Info("----------IKuuu签到结果：" + msg + "----------")
	logger.Log().Info("----------IKuuu结束签到----------")
	return msg
}
