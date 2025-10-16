package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/notifier"
	"auto-checkin/internal/scheduler"
)

func main() {
	// 加载配置
	_, err := config.Init("config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志模块
	err = logger.Log().Init("logs/app.log")
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	logger.Log().Info("服务已启动")
	defer logger.Log().Close()
	// 初始化推送模块
	notify := notifier.New()
	// 初始化定时任务
	sd := scheduler.New(notify)
	sd.Start()

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 停止定时任务
	sd.Stop()
	logger.Log().Info("服务已停止")
}
