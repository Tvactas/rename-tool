package applog

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// logFileWriter 封装日志文件写入
type logFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// appdata -> user -> D:/
func newLogFileWriter(fileName string) *logFileWriter {
	lw := &logFileWriter{}

	// 1. 优先使用 Windows 用户 AppData 目录（规范做法）
	userDir := getUserDir()
	path := filepath.Join(userDir, fileName)
	if file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
		lw.file = file
		lw.writer = bufio.NewWriter(file)
		return lw
	}

	// 2. D盘
	path = filepath.Join("D:\\", fileName)
	if file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
		lw.file = file
		lw.writer = bufio.NewWriter(file)
		return lw
	}

	// 3. 都失败则丢弃
	lw.writer = bufio.NewWriter(io.Discard)
	return lws
}

func getUserDir() string {
	// 先尝试获取 APPDATA
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return appData
	}

	// 如果 APPDATA 不存在，则退回用户主目录
	home, err := os.UserHomeDir()
	if err != nil {
		// 最坏情况：返回当前目录
		return "."
	}
	return home
}
