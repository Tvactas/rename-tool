package menu

import (
	"rename-tool/common/applog"
	"rename-tool/setting/i18n"
)

// -------------------- 国际化 & 日志 --------------------

func buttonTr(key string) string {
	return i18n.ButtonTr(key)
}

func logTr(key string) string {
	return i18n.LogTr(key)
}

func logEvent(prefix, key string) {
	applog.Logger.Printf("[%s] %s", prefix, logTr(key))
}
