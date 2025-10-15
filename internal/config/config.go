package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Website struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
	Cookies map[string]string `yaml:"cookies"`
}

type WeCom struct {
	Webhook string `yaml:"webhook"`
}

type Telegram struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

type Notifications struct {
	WeCom    WeCom    `yaml:"wecom"`
	Telegram Telegram `yaml:"telegram"`
}

type Config struct {
	Websites      []Website     `yaml:"websites"`
	Notifications Notifications `yaml:"notifications"`
}

var Cfg = &Config{}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, Cfg); err != nil {
		return nil, err
	}
	return Cfg, nil
}
