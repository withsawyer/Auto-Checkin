package handler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"encoding/json"
	"fmt"
)

func init() {
	RegisterCheckInHandler("quark", &Quark{}) // 注册处理器
}

// Quark 封装夸克签到逻辑
type Quark struct {
	BaseLogic
	website cfg.Website
	//Config QuarkConfig // 夸克网盘配置信息
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
		Headers:            q.website.Headers,
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
		Method:             "GET",
		URL:                "https://drive-m.quark.cn/1/clouddrive/capacity/growth/info",
		QueryParams:        q.website.Query,
		Headers:            q.website.Headers,
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
	jsonData, err := json.Marshal(q.website.Body)
	if err != nil {
		return false, "", err
	}
	response, err := util.SendRequest(&util.RequestParams{
		Method:             "POST",
		URL:                "https://drive-m.quark.cn/1/clouddrive/capacity/growth/sign",
		QueryParams:        q.website.Query,
		BodyData:           jsonData,
		Headers:            q.website.Headers,
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
func (q *Quark) doSign() error {
	userinfo := q.getUserInfo()
	if userinfo == nil {
		q.PushContent("❌ 获取用户信息失败")
	}
	// 获取签到信息
	growthInfo, err := q.getGrowthInfo()
	if err != nil {
		q.PushContent("❌ 获取成长信息失败")
		return fmt.Errorf("❌ 获取成长信息失败: %v", err)
	}
	// 记录用户信息
	isVIP := "普通用户"
	if growthInfo["88VIP"].(bool) {
		isVIP = "88VIP"
	} else if growthInfo["super_vip_exp_at"].(float64) > 0 {
		isVIP = "SVIP"
	}
	// 昵称兼容显示
	nickname := userinfo["nickname"].(string)
	if nickname == "" {
		nickname = "查询失败"
	}
	q.PushContent("👶 用户名: %s[%s]", nickname, isVIP)
	// 记录容量信息
	totalCapacity := growthInfo["total_capacity"].(float64)
	q.PushContent("💾 网盘总容量: %s，", q.convertBytes(int64(totalCapacity)))

	if capComp, ok := growthInfo["cap_composition"].(map[string]interface{}); ok {
		if reward, ok := capComp["sign_reward"].(float64); ok {
			q.PushContent("✏️ 签到累计容量: %s", q.convertBytes(int64(reward)))
		} else {
			q.PushContent("✏️ 签到累计容量: 0 MB")
		}
	}
	// 检查是否已签到
	if capSign, ok := growthInfo["cap_sign"].(map[string]interface{}); ok {
		if signed := capSign["sign_daily"].(bool); signed {
			reward := capSign["sign_daily_reward"].(float64)
			progress := capSign["sign_progress"].(float64)
			target := capSign["sign_target"].(float64)
			q.PushContent("✅ 签到日志: 今日已签到+%s，连签进度(%.0f/%.0f)", q.convertBytes(int64(reward)), progress, target)
		} else {
			success, reward, err := q.getGrowthSign()
			if err != nil {
				q.PushContent("❌ 签到异常")
				logger.Log().Errorf("❌ 签到异常: %v", err)
			} else if success {
				progress := capSign["sign_progress"].(float64) + 1
				target := capSign["sign_target"].(float64)
				q.PushContent("✅ 执行签到: 今日签到+%s，连签进度(%.0f/%.0f)", reward, progress, target)
			} else {
				q.PushContent("❌ 签到异常")
				logger.Log().Errorf("❌ 签到异常: %s\n", reward)
			}
		}
	}
	return nil
}

// NewQuark 初始化 Quark 实例
func NewQuark(website cfg.Website) *Quark {
	obj := &Quark{
		website: website,
	}
	obj.Content = "👙 [服务]" + website.Name + "签到信息\n"
	return obj
}

func (q *Quark) Run(website cfg.Website) string {
	logger.Log().Debug("----------夸克网盘开始签到----------")
	// 执行签到
	quark := NewQuark(website)
	res := quark.doSign()
	if res != nil {
		logger.Log().Error("签到失败: " + res.Error())
	}
	logger.Log().Debug("----------夸克网盘签到完毕----------")
	return quark.Content
}
