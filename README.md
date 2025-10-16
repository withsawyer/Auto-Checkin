# Auto-Checkin

一个简单的自动签到脚本，用于帮助用户自动完成日常签到任务。支持多种网站和平台，通过配置文件灵活管理任务。

## 功能

- 支持多平台签到（如京东、Quark、iKuuu等）。
- 通过配置文件动态加载任务。
- 支持定时任务调度。
- 提供日志记录和通知功能（如企业微信、Telegram）。

## 项目结构

```
├── config.json            # 主配置文件
├── config.json.example    # 配置文件示例
├── go.mod                 # Go模块定义
├── go.sum                 # 依赖版本锁定
├── main.go                # 程序入口
├── README.md              # 项目说明
└── internal/              # 内部模块
    ├── config/            # 配置管理
    ├── handler/           # 签到处理器
    ├── interfaces/        # 接口定义
    ├── logger/            # 日志管理
    ├── notifier/          # 通知管理
    ├── scheduler/         # 任务调度
    └── util/              # 工具函数
```

## 快速开始

1. 复制 `config.json.example` 为 `config.json`，并根据需要修改配置。
2. 运行 `go run main.go` 启动程序。

## 配置说明

配置文件 `config.json` 包含以下主要部分：

- `websites`: 定义需要签到的网站信息（如请求头、参数、Cookie等）。
- `notifiers`: 配置通知方式（如企业微信、Telegram）。
- `cron`: 定义定时任务规则。

## 示例配置

```json
{
  "websites": [
    {
      "name": "jd",
      "headers": {
        "User-Agent": "Mozilla/5.0"
      },
      "cookies": {
        "key": "value"
      }
    }
  ],
  "notifiers": {
    "wecom": {
      "webhook": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"
    }
  },
  "cron": "0 9 * * *"
}
```

## 开发指南

1. **添加新平台**：在 `internal/handler/` 下实现新的签到处理器，并注册到 `init` 函数中。
2. **扩展通知方式**：在 `internal/notifier/` 下实现新的通知逻辑。
3. **调试**：使用 `logger` 模块记录日志，便于排查问题。

## 依赖

- Go 1.16+
- 第三方库：`github.com/robfig/cron/v3`（定时任务）

## 许可证

MIT