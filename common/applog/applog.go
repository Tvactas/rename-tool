package applog

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"rename-tool/setting/config"
	"sync"
)

var (
	// Logger 全局日志实例
	Logger *log.Logger
	once   sync.Once

	logWriter *logFileWriter
)

// logFileWriter 封装日志文件写入
type logFileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// newLogFileWriter 创建日志写入对象（按优先级）
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
	return lw
}

// Write 实现 io.Writer 接口
func (lw *logFileWriter) Write(p []byte) (n int, err error) {
	if lw.writer != nil {
		n, err = lw.writer.Write(p)
		if err == nil {
			lw.writer.Flush()
		}
		return n, err
	}
	return len(p), nil
}

// Close 刷新缓存并关闭文件
func (lw *logFileWriter) Close() {
	if lw.writer != nil {
		lw.writer.Flush()
	}
	if lw.file != nil {
		lw.file.Close()
	}
}

// InitLogger 初始化全局 logger，只执行一次
func InitLogger(fileName string) {
	once.Do(func() {
		logWriter = newLogFileWriter(fileName)
		Logger = log.New(logWriter, "", log.LstdFlags)

		// 程序退出时自动关闭
		go func() {
			c := make(chan os.Signal, 1)
			<-c // 可自行处理退出信号
			logWriter.Close()
		}()
	})
}

// GetUserDir 返回程序的用户数据目录（Windows 专用）。
// 优先使用 APPDATA，如果没有则回退到用户主目录。
// 最终路径形如：C:\Users\<user>\AppData\Roaming\<AppName>
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

func GetLogPath() string {
	appDir := getUserDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "rename.log")
	}
	return filepath.Join(logDir, "rename.log")
}
