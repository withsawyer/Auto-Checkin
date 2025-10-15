package scheduler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"encoding/json"
	"fmt"
)

// Config 定义配置文件结构
type Config struct {
	User    string            `json:"user"`    // 用户名称
	URL     string            `json:"url"`     // 夸克网盘URL
	KPS     string            `json:"kps"`     // KPS参数
	Sign    string            `json:"sign"`    // 签名参数
	VCode   string            `json:"vcode"`   // 验证码参数
	Headers map[string]string `json:"headers"` // 自定义HTTP Header
}

func init() {
	RegisterCheckInHandler("quark", &Quark{}) // 注册处理器
}

// Quark 封装夸克签到逻辑
type Quark struct {
	Config Config // 夸克网盘配置信息
}

// convertBytes 将字节转换为 MB/GB/TB
func (q *Quark) convertBytes(b int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"} // 单位列表
	i := 0
	for b >= 1024 && i < len(units)-1 {
		b /= 1024 // 转换为更高一级单位
		i++
	}
	return fmt.Sprintf("%.2f %s", float64(b), units[i]) // 返回格式化后的字符串
}

// getUserInfo 获取用户信息
func (q *Quark) getUserInfo() map[string]interface{} {
	result, err := util.SendRequest(&util.RequestParams{
		Method: "GET",
		URL:    "https://pan.quark.cn/account/info",
		QueryParams: map[string]string{
			"fr":       "pc",
			"platform": "pc",
		},
		Headers:            q.Config.Headers,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		return data
	}
	return nil
}

// getGrowthInfo 获取用户当前的签到信息
func (q *Quark) getGrowthInfo() (map[string]interface{}, error) {
	result, err := util.SendRequest(&util.RequestParams{
		Method: "GET",
		URL:    "https://drive-m.quark.cn/1/clouddrive/capacity/growth/info",
		QueryParams: map[string]string{
			"pr":    "ucpro",
			"fr":    "android",
			"kps":   q.Config.KPS,
			"sign":  q.Config.Sign,
			"vcode": q.Config.VCode,
		},
		Headers:            q.Config.Headers,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return nil, fmt.Errorf("failed to get growth info")
}

// getGrowthSign 执行签到
func (q *Quark) getGrowthSign() (bool, string, error) {
	jsonData, err := json.Marshal(map[string]string{
		"pr":    "ucpro",
		"fr":    "android",
		"kps":   q.Config.KPS,
		"sign":  q.Config.Sign,
		"vcode": q.Config.VCode,
	})
	if err != nil {
		return false, "", err
	}
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                "https://drive-m.quark.cn/1/clouddrive/capacity/growth/sign",
		BodyData:           jsonData,
		Headers:            q.Config.Headers,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return false, "", err
	}

	if data, ok := response["data"].(map[string]interface{}); ok {
		reward := data["sign_daily_reward"].(float64)
		return true, q.convertBytes(int64(reward)), nil
	}
	return false, response["message"].(string), nil
}

// DoSign 执行签到任务
func (q *Quark) doSign() string {
	message := ""

	userinfo := q.getUserInfo()
	if userinfo == nil {
		message += fmt.Sprintf("❌ 获取用户信息失败\n")
		return message
	}
	message += fmt.Sprintf("🔔 用户名: %s\n", userinfo["nickname"].(string))

	// 获取签到信息
	growthInfo, err := q.getGrowthInfo()
	if err != nil {
		message += fmt.Sprintf("❌ 获取成长信息失败: %v\n", err)
		return message
	}

	// 记录用户信息
	isVIP := "普通用户"
	if growthInfo["88VIP"].(bool) {
		isVIP = "88VIP"
	} else if growthInfo["super_vip_exp_at"].(float64) > 0 {
		isVIP = "SVIP"
	}
	message += fmt.Sprintf("%s\n", isVIP)

	// 记录容量信息
	totalCapacity := growthInfo["total_capacity"].(float64)
	message += fmt.Sprintf("💾 网盘总容量: %s，", q.convertBytes(int64(totalCapacity)))

	if capComp, ok := growthInfo["cap_composition"].(map[string]interface{}); ok {
		if reward, ok := capComp["sign_reward"].(float64); ok {
			message += fmt.Sprintf("签到累计容量: %s\n", q.convertBytes(int64(reward)))
		} else {
			message += "签到累计容量: 0 MB\n"
		}
	}
	// 检查是否已签到
	if capSign, ok := growthInfo["cap_sign"].(map[string]interface{}); ok {
		if signed := capSign["sign_daily"].(bool); signed {
			reward := capSign["sign_daily_reward"].(float64)
			progress := capSign["sign_progress"].(float64)
			target := capSign["sign_target"].(float64)
			message += fmt.Sprintf("✅ 签到日志: 今日已签到+%s，连签进度(%.0f/%.0f)\n",
				q.convertBytes(int64(reward)), progress, target)
		} else {
			success, reward, err := q.getGrowthSign()
			if err != nil {
				message += fmt.Sprintf("❌ 签到异常: %v\n", err)
			} else if success {
				progress := capSign["sign_progress"].(float64) + 1
				target := capSign["sign_target"].(float64)
				message += fmt.Sprintf("✅ 执行签到: 今日签到+%s，连签进度(%.0f/%.0f)\n", reward, progress, target)
			} else {
				message += fmt.Sprintf("❌ 签到异常: %s\n", reward)
			}
		}
	}

	return message
}

// NewQuark 初始化 Quark 实例
func NewQuark(website cfg.Website) *Quark {
	var config Config
	err := json.Unmarshal([]byte(website.Body), &config)
	if err != nil {
		return nil
	}
	config.URL = website.URL
	config.Headers = website.Headers
	return &Quark{
		Config: config,
	}
}

func (q *Quark) Run(website cfg.Website) string {
	logger.Log().Info("----------夸克网盘开始签到----------")
	// 执行签到
	quark := NewQuark(website)
	msg := quark.doSign()
	logger.Log().Info(msg)
	logger.Log().Info("----------夸克网盘签到完毕----------")
	return msg
}
