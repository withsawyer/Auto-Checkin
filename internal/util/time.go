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
