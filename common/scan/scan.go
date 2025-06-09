package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"rename-tool/common/FileStatus"
	"rename-tool/common/log"
	"sort"
	"strings"
)

func ScanFormats(dir string) ([]string, error) {
	formatMap := make(map[string]struct{})
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 尝试打开文件以确保可访问
			file, err := os.Open(path)
			if err != nil {
				// 如果文件被占用，记录错误但继续处理其他文件
				if FileStatus.IsFileBusyError(err) {
					log.LogError(fmt.Errorf("file busy: %s", path))
					return nil
				}
				return err
			}
			file.Close()

			ext := strings.ToLower(filepath.Ext(path))
			if ext != "" {
				formatMap[ext] = struct{}{}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	formats := make([]string, 0, len(formatMap))
	for ext := range formatMap {
		formats = append(formats, ext)
	}
	sort.Strings(formats)
	return formats, nil
}
