package filestatus

import (
	"errors"
	"syscall"
)

const errorSharingViolation syscall.Errno = 32 // Windows ERROR_SHARING_VIOLATION

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
