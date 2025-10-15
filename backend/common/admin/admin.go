package admin

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	initOnce sync.Once
	logger   *log.Logger
)

// Initialize logger
func init() {
	logOutput := getLogWriter()
	logger = log.New(logOutput, "", log.LstdFlags)
}

// getLogWriter returns a writer to a log file with fallback
func getLogWriter() io.Writer {
	paths := []string{}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, "tvacats_rename.log"))
	}
	paths = append(paths, "D:\\tvacats_rename.log")

	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			return file
		}
	}

	// Fallback: write to stderr
	return os.Stderr
}

// Log logs a formatted message using the package logger
func Log(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// IsAdmin returns true if running as administrator, false otherwise
func IsAdmin() bool {
	var isAdmin bool
	initOnce.Do(func() {
		// Create SID for Administrators group
		sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
		if err != nil {
			Log("CreateWellKnownSid failed: %v", err)
			return
		}

		// Get current process token
		token := windows.GetCurrentProcessToken() // 非 deprecated
		defer token.Close()

		isMember, err := token.IsMember(sid)
		if err != nil {
			Log("IsMember failed: %v", err)
			return
		}

		isAdmin = isMember
		Log("Administrator check result: %v", isAdmin)
	})
	return isAdmin
}

// Example prints the current admin status
func Example() {
	if IsAdmin() {
		println("Running as administrator")
	} else {
		println("Running as normal user")
	}
}
