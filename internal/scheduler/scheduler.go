package scheduler

import (
	"auto-checkin/internal/config"
	"auto-checkin/internal/handler"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"github.com/robfig/cron/v3"
	"strings"
	"sync"
	"time"

	"auto-checkin/internal/notifier"
)

type Scheduler struct {
	notifier *notifier.Notifier
	ticker   *time.Ticker
	done     chan bool
}

func New(notifier *notifier.Notifier) *Scheduler {
	return &Scheduler{
		notifier: notifier,
	}
}

func (s *Scheduler) Start() {
	if config.Cfg.Debug {
		// 调试模式只执行一次
		s.runCheckIn()
	} else {
		c := cron.New(cron.WithLocation(util.GetTimeLocation()))
		_, err := c.AddFunc(config.Cfg.Cron, s.runCheckIn)
		if err != nil {
			logger.Log().Error("定时任务配置错误: " + err.Error())
			return
		}
		c.Start()
		logger.Log().Info("定时任务已启动，执行规则: " + config.Cfg.Cron)
		select {}
	}
}

func (s *Scheduler) runCheckIn() {
	logger.Log().Info("开始签到任务")
	var wg sync.WaitGroup
	var handlers []string
	for h, _ := range handler.CheckinHandlers {
		handlers = append(handlers, h)
	}
	logger.Log().Debugf("当前注册的处理器: %+v", handlers)

	var signRes []string
	for index, website := range config.Cfg.Websites {
		wg.Add(1)
		go func(i int, w config.Website) {
			defer wg.Done()
			handle, ok := handler.CheckinHandlers[strings.ToLower(w.Name)]
			if !ok {
				signRes = append(signRes, "\n[服务] "+w.Name+"\n❌ 不支持的签到服务: "+w.Name+"\n")
				logger.Log().Info("不支持的签到服务: " + w.Name)
			} else {
				logger.Log().Info("开始签到: " + w.Name)
				m := handle.Run(w)
				if m != "" {
					signRes = append(signRes, "\n"+m)
				}
				logger.Log().Info("签到完成: " + w.Name)
			}

		}(index, website)
	}
	signContent := "\n≡≡≡≡≡≡ 签到任务报告 ≡≡≡≡≡≡\n"
	wg.Wait()
	for i, re := range signRes {
		signContent += re
		if i < len(signRes)-1 {
			signContent += "\n—————————————\n"
		}
	}
	signContent += "\n≡≡≡≡≡≡ 任务结束 ≡≡≡≡≡≡"
	logger.Log().Debug(signContent)
	s.notifier.Push(signContent)
}
