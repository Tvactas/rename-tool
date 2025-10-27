package dirpath

import (
	"strings"
)

// truncatePathMiddle 截断路径，保留前后信息，中间用省略号
func truncatePathMiddle(path string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}

	runes := []rune(path) // 支持多字节字符
	n := len(runes)

	if n <= maxLength {
		return path
	}

	if maxLength <= 3 {
		return strings.Repeat(".", maxLength)
	}

	frontLen := (maxLength - 3) / 2
	backLen := maxLength - 3 - frontLen

	return string(runes[:frontLen]) + "..." + string(runes[n-backLen:])
}

// map文件夹下的所有格式
func mapExt(formats []string) map[string]bool {
	m := make(map[string]bool, len(formats))
	for _, ext := range formats {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		m[ext] = true
	}
	return m
}
