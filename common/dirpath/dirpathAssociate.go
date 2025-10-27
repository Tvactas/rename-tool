package dirpath

import (
	"rename-tool/common/applog"
	"rename-tool/setting/i18n"
)

func logTr(key string) string {
	return i18n.LogTr(key)
}

func logEvent(prefix, key string, value any) {
	applog.Logger.Printf("[%s] %s: %v", prefix, logTr(key), value)
}
