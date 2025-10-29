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

	userDir := getUserDir()
	path := filepath.Join(userDir, fileName)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err == nil {
		lw.file = file
		lw.writer = bufio.NewWriter(file)
		return lw
	}

	// 如果打开失败，直接使用 io.Discard 避免程序崩溃
	lw.writer = bufio.NewWriter(io.Discard)
	return lw
}

func getUserDir() string {
	// 先尝试 APPDATA
	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData
	}

	// 尝试用户主目录
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}

	// 最后退回当前目录
	return "."
}
