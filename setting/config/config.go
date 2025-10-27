package config

import "time"

const (
	LogDir = "logs"

	// 文件操作相关常量
	MaxRetryAttempts = 3
	RetryDelay       = 500 * time.Millisecond

	// UI 相关常量
	DefaultWindowWidth  = 600
	DefaultWindowHeight = 400
	ProgressBarWidth    = 300
	ButtonMinWidth      = 120
	ButtonMinHeight     = 28
	DialogMinWidth      = 400
	DialogMinHeight     = 300
	FormatListHeight    = 200
)
