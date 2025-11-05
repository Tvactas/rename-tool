package dirpath

import (
	"fmt"
	"os"
	"path/filepath"
	"rename-tool/common/filestatus"
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

// walkDirFiltered 统一封装遍历逻辑：
// 按扩展名过滤文件，并为每个文件调用 fn(name string)。
func walkDirFilteredWalk(root string, formats []string, fn func(path string, info os.FileInfo)) error {
	formatsMap := mapExt(formats)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if filestatus.IsFileBusyError(err) {
				return nil // 忽略占用错误
			}
			return fmt.Errorf("%s: %w", textTr("failReadFiles"), err)
		}
		if info.IsDir() {
			return nil
		}
		if len(formatsMap) == 0 || formatsMap[strings.ToLower(filepath.Ext(path))] {
			fn(path, info)
		}
		return nil
	})
}

func walkDirFiltered(root string, formats []string, fn func(path string, info os.FileInfo)) error {
	formatsMap := mapExt(formats)

	entries, err := os.ReadDir(root)
	if err != nil {
		if filestatus.IsFileBusyError(err) {
			return nil // 忽略占用错误
		}
		return fmt.Errorf("%s: %w", textTr("failReadFiles"), err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			if filestatus.IsFileBusyError(err) {
				continue
			}
			return fmt.Errorf("%s: %w", textTr("failReadFiles"), err)
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if len(formatsMap) == 0 || formatsMap[ext] {
			fn(filepath.Join(root, entry.Name()), info)
		}
	}

	return nil
}
