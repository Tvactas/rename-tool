package scan

import (
	"os"
	"path/filepath"
	"rename-tool/common/applog"
	"rename-tool/common/filestatus"
	"rename-tool/setting/i18n"
	"sort"
	"strings"
)

func ScanFormats(dir string) ([]string, error) {
	formatMap := make(map[string]struct{})
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if filestatus.IsFileBusyError(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			// 尝试打开文件以确保可访问
			file, err := os.Open(path)
			if err != nil {
				// 如果文件被占用，记录错误但继续处理其他文件
				if filestatus.IsFileBusyError(err) {
					applog.Logger.Printf("[FILE ERROR] %s,%v", i18n.LogTr("FileStatus"), path)
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
