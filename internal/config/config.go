package config

import (
	"encoding/json"
	"os"
)

type Website struct {
	Name    string            `json:"name"`
	Headers map[string]string `json:"headers"`
	Query   map[string]string `json:"query"`
	Body    map[string]any    `json:"body"`
	Cookies map[string]string `json:"cookies"`
}

type WeCom struct {
	KEY string `json:"key"`
}

type Telegram struct {
	BotToken string `json:"bot_token"`
	UID      string `json:"uid"`
	APIHost  string `json:"api_host"`
	ChatID   string `json:"chat_id"`
}

type Notifications struct {
	WeCom    WeCom    `json:"wecom"`
	Telegram Telegram `json:"telegram"`
}
type Proxy struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
type Config struct {
	Cron          string        `json:"cron"`
	Debug         bool          `json:"debug"`
	Websites      []Website     `json:"websites"`
	Notifications Notifications `json:"notifications"`
	Proxy         Proxy         `json:"proxy"`
}

var Cfg = &Config{}

func Init(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, Cfg); err != nil {
		return nil, err
	}
	return Cfg, nil
}
