package dirpath

import (
	"fmt"
	"os"
	"path/filepath"

	"rename-tool/common/FileStatus"

	"rename-tool/setting/config"
	"sort"
	"strings"
)

// 获取应用数据目录
func GetAppDataDir() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	// 在用户主目录下创建应用目录
	appDir := filepath.Join(homeDir, "."+config.AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "."
	}
	return appDir
}

func GetFiles(dir string, formats []string) ([]string, error) {
	var result []string
	formatSet := make(map[string]bool)
	for _, f := range formats {
		formatSet[f] = true
	}

	// 使用缓冲通道优化文件遍历
	fileChan := make(chan string, 100)
	errorChan := make(chan error, 1)

	go func() {
		defer close(fileChan)
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
						return fmt.Errorf("file busy: %s", path)
					}
					return err
				}
				file.Close()

				fileChan <- path
			}
			return nil
		})
		if err != nil {
			errorChan <- err
		}
	}()

	// 处理文件
	for file := range fileChan {
		ext := strings.ToLower(filepath.Ext(file))
		if len(formats) == 0 || formatSet[ext] {
			result = append(result, file)
		}
	}

	// 检查错误
	select {
	case err := <-errorChan:
		return nil, err
	default:
	}

	sort.Strings(result)
	return result, nil
}
