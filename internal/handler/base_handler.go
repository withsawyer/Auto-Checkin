package handler

import (
	"auto-checkin/internal/interfaces"
	"fmt"
	"strings"
)

type BaseLogic struct {
	Content string
}

func (b *BaseLogic) PushContent(format string, args ...any) string {
	content := fmt.Sprintf(format, args...)
	b.Content += "∷∷∷∷" + content + "\n"
	return b.Content
}

// CheckinHandlers  全局工厂，存储所有签到处理器
var CheckinHandlers = make(map[string]interfaces.Logic)

// RegisterCheckInHandler 注册签到处理器
func RegisterCheckInHandler(name string, handler interfaces.Logic) {
	CheckinHandlers[strings.ToLower(name)] = handler
}
