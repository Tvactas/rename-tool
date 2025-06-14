package fs

import (
	"errors"
	"strings"
	"syscall"
)

// AppError represents an application error
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// IsFileBusyError checks if the error is a file busy error
func IsFileBusyError(err error) bool {
	if err == nil {
		return false
	}

	// Check Windows system error
	if errors.Is(err, syscall.Errno(32)) { // ERROR_SHARING_VIOLATION = 32
		return true
	}

	// Check error message for keywords
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "process cannot access the file") ||
		strings.Contains(errMsg, "file is being used by another process") ||
		strings.Contains(errMsg, "sharing violation") ||
		strings.Contains(errMsg, "file is locked")
}
