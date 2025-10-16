package util

import (
	"auto-checkin/internal/logger"
	"time"
)

func GetTimeLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Log().Error("时区加载失败: " + err.Error())
		return loc
	}
	return time.Local
}

func GetNowUnixTimestamp() int64 {
	currentTime := time.Now()
	// 获取精确到秒的时间戳
	timestamp := currentTime.Unix()
	return timestamp
}

func GetMilliTimestamp() int64 {
	currentTime := time.Now()
	// 获取纳秒时间戳
	nanoTimestamp := currentTime.UnixNano()
	// 将纳秒时间戳转换为毫秒时间戳
	milliTimestamp := nanoTimestamp / 1_000_000
	return milliTimestamp
}
