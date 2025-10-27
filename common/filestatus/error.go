package filestatus

import (
	"errors"
	"strings"
	"syscall"
)

const errorSharingViolation syscall.Errno = 32 // Windows ERROR_SHARING_VIOLATION

// AppError represents an application-level error with code and message.
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

// Unwrap allows AppError to participate in error chains.
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsFileBusyError checks whether the error indicates a "file is in use" condition.
func IsFileBusyError(err error) bool {
	if err == nil {
		return false
	}

	// Check known Windows error code
	if errors.Is(err, errorSharingViolation) {
		return true
	}

	// Fallback: fuzzy matching of common error message substrings
	return isKnownFileBusyMessage(err.Error())
}

// isKnownFileBusyMessage checks if the message contains known "file busy" patterns.
func isKnownFileBusyMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "process cannot access the file") ||
		strings.Contains(msg, "file is being used by another process") ||
		strings.Contains(msg, "sharing violation") ||
		strings.Contains(msg, "file is locked")
}
