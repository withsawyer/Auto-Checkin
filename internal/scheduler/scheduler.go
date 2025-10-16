package scheduler

import (
	"auto-checkin/internal/config"
	"auto-checkin/internal/handler"
	"auto-checkin/internal/logger"
	"auto-checkin/internal/util"
	"fmt"
	"github.com/robfig/cron/v3"
	"regexp"
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
		done:     make(chan bool),
	}
}

func (s *Scheduler) Start() {
	if config.Cfg.Debug {
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
	}
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
}

func (s *Scheduler) runCheckIn() {
	signContent := "\n========== 签到任务报告 ==========\n"
	logger.Log().Info("开始签到任务")
	var wg sync.WaitGroup
	logger.Log().Debug("当前注册的处理器: %+v", handler.CheckinHandlers)
	fmt.Printf("%+v\n", handler.CheckinHandlers)
	for index, website := range config.Cfg.Websites {
		wg.Add(1)
		go func(i int, w config.Website) {
			defer wg.Done()
			handle, ok := handler.CheckinHandlers[strings.ToLower(w.Name)]
			if !ok {
				signContent += "\n[服务] " + w.Name + "\n❌ 不支持的签到服务: " + w.Name + "\n"
				logger.Log().Info("不支持的签到服务: " + w.Name)
			} else {
				logger.Log().Info("开始签到: " + w.Name)
				m := handle.Run(w)
				if m != "" {
					signContent += "\n" + m
					if index < len(config.Cfg.Websites)-1 {
						signContent += "\n\n☆☆☆☆☆☆☆☆☆☆☆☆☆☆☆\n"
					}
				}
				logger.Log().Info("签到完成: " + w.Name)
			}

		}(index, website)
	}
	wg.Wait()
	signContent += "\n========== 任务结束 =========="
	logger.Log().Debug(signContent)
	_ = s.notifier.SendTelegram(signContent)
	_ = s.notifier.SendWeCom(signContent)
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
