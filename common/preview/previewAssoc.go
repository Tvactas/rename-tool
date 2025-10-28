package preview

import "rename-tool/setting/i18n"

// tr 函数用于国际化
func tr(key string) string {
	return i18n.Tr(key)
}
