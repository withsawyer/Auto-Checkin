package scheduler

import (
	"auto-checkin/internal/config"
	"auto-checkin/internal/logger"
	"fmt"
	"regexp"
	"time"

	"auto-checkin/internal/http"
	"auto-checkin/internal/notifier"
)

type Scheduler struct {
	websites []config.Website
	client   *http.Client
	notifier *notifier.Notifier
	ticker   *time.Ticker
	done     chan bool
}

func New(websites []config.Website, client *http.Client, notifier *notifier.Notifier) *Scheduler {
	return &Scheduler{
		websites: websites,
		client:   client,
		notifier: notifier,
		done:     make(chan bool),
	}
}

func (s *Scheduler) Start() {
	//s.ticker = time.NewTicker(24 * time.Hour)
	//s.ticker = time.NewTicker(1 * time.Second)
	s.runCheckIn()
	fmt.Printf("定时任务已启动\n")
	//go func() {
	//	for {
	//		select {
	//		case <-s.ticker.C:
	//			print("定时任务已执行")
	//			s.runCheckIn()
	//		case <-s.done:
	//			return
	//		}
	//	}
	//}()
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
}

// CheckInHandler 定义签到处理器接口
type CheckInHandler interface {
	Run(website config.Website) string
}

// checkInHandlers 全局工厂，存储所有签到处理器
var checkInHandlers = make(map[string]CheckInHandler)

// RegisterCheckInHandler 注册签到处理器
func RegisterCheckInHandler(name string, handler CheckInHandler) {
	checkInHandlers[name] = handler
}

func (s *Scheduler) runCheckIn() {
	msg := "签到结果:"
	for _, website := range s.websites {
		msg += "\n　"
		go func(w config.Website) {
			handler, ok := checkInHandlers[w.Name]
			if !ok {
				logger.Log().Info("不支持的签到服务: " + w.Name)
				return
			}
			logger.Log().Info("开始签到: " + w.Name)
			m := handler.Run(w)
			if m != "" {
				msg += m
			}
			logger.Log().Info("签到完成: " + w.Name)
		}(website)
	}
	fmt.Println(msg)
	logger.Log().Info(msg)
}

func (s *Scheduler) matchLogic(url string, match string) bool {
	// 创建正则表达式模式
	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(match))
	// 检查 URL 是否以指定前缀开头
	if re.MatchString(url) {
		return true
	} else {
		return false
	}
}
