package applog

import (
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

		go func() {
			c := make(chan os.Signal, 1)
			<-c
			logWriter.Close()
		}()
	})
}

func GetLogPath() string {
	appDir := getUserDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "rename.log")
	}
	return filepath.Join(logDir, "rename.log")
}
