package interfaces

import "auto-checkin/internal/config"

type Logic interface {
	Run(website config.Website) string
	PushContent(format string, args ...any) string
}
