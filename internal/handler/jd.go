package handler

import (
	cfg "auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"fmt"
	"net/url"
)

func init() {
	RegisterCheckInHandler("jd", &JD{}) // 注册处理器
}

// JD 封装京东签到逻辑
type JD struct {
	BaseLogic
	website cfg.Website
}

func (j *JD) balance() error {
	reqData := url.Values{}
	reqData.Add("appid", j.website.Body["appid"].(string))
	reqData.Add("functionId", "BEAN_BALANCE")
	reqData.Add("body", "{}")
	reqData.Add("client", j.website.Body["client"].(string))
	reqData.Add("_t", "1760671276136")

	reqParams := &util.RequestParams{
		Method: "POST",
		URL:    "https://api.m.jd.com/api",
		QueryParams: map[string]string{
			"functionId": "BEAN_BALANCE",
			"appid":      "asset-h5",
		},
		BodyData:           reqData,
		Headers:            j.website.Headers,
		InsecureSkipVerify: true,
	}
	logger.Log().Debug("⌛ 准备发送BEAN_BALANCE请求")
	response, err := util.SendRequest(reqParams)
	if err != nil {
		logger.Log().Error("❌ 发送请求失败")
		return fmt.Errorf("❌ 发送请求失败: %v", err)
	}
	if code, ok := response["code"].(string); ok && "0000" == code {
		if data, ok := response["data"].(map[string]any); ok && data != nil {
			j.PushContent("🍅 京豆余额:%.f", data["balance"].(float64))
			return nil
		}
	}
	j.PushContent("🍅 京豆余额:获取失败")
	return fmt.Errorf("❌ 获取京豆余额失败: %v", response)
}

// doSign 执行京东签到任务
func (j *JD) doSign() error {
	// 构造请求参数
	values, err := util.Map2UrlValues(j.website.Body)
	if err != nil {
		j.PushContent("❌ [doSign]构造请求参数失败")
		return err
	}

	reqParams := &util.RequestParams{
		Method:             "POST",
		URL:                "https://api.m.jd.com/",
		BodyData:           values,
		Headers:            j.website.Headers,
		InsecureSkipVerify: true,
	}
	response, err := util.SendRequest(reqParams)
	if err != nil {
		j.PushContent("❌ 发送请求失败")
		return fmt.Errorf("❌ 发送请求失败: %v", err)
	}

	// 处理签到结果
	if success, ok := response["success"].(bool); ok && success == true {
		responseData := response["data"].(map[string]any)
		if assignmentInfo, ok := responseData["assignmentInfo"].(map[string]any); ok && assignmentInfo != nil {
			j.PushContent("📋 累计签到次数:%d", assignmentInfo["completionCnt"].(int))
			j.PushContent("📋 连续签到次数:%d", assignmentInfo["continueSignDay"].(int))
		}
		if assignmentRewardInfo, ok := responseData["assignmentRewardInfo"].(map[string]any); ok && assignmentRewardInfo != nil {
			if jingDouRewards, ok := responseData["jingDouRewards"].([]map[string]any); ok && jingDouRewards != nil {
				for _, reward := range jingDouRewards {
					j.PushContent("🏆 签到奖励:%s", reward["rewardName"].(string))
				}
			}
		}
		j.PushContent("✅ 京东签到成功")
		return nil
	} else {
		if errCode, ok := response["errCode"].(string); ok && errCode == "302" {
			j.PushContent("✅ 京东已完成签到")
			return nil
		} else {
			if errMessage, ok := response["errMessage"].(string); ok {
				j.PushContent("📞 %s", errMessage)
			}
		}
	}
	j.PushContent("❌ 京东签到失败")
	return fmt.Errorf("❌ 京东签到失败: %v", response["message"])
}

// NewJD 初始化 JD 实例
func NewJD(website cfg.Website) *JD {
	website.Body["t"] = util.GetMilliTimestamp()
	obj := &JD{
		website: website,
	}
	obj.Content = "👙 [服务]" + website.Name + "签到信息\n"
	return obj
}

func (j *JD) Run(website cfg.Website) string {
	logger.Log().Debug("----------京东开始签到----------")
	// 执行签到
	jd := NewJD(website)
	if jd == nil {
		logger.Log().Error("❌ 京东初始化失败")
		return "❌ 京东初始化失败"
	}
	_ = jd.balance()
	res := jd.doSign()
	if res != nil {
		logger.Log().Error("签到失败: " + res.Error())
	}
	return jd.Content
}
