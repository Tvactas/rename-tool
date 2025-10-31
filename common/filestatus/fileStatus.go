package filestatus

import (
	"errors"
	"fmt"
	"os"
	"rename-tool/common/antisamename"
	"rename-tool/setting/config"
	"rename-tool/setting/i18n"
	"syscall"
	"time"
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

func RenameFile(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}
	newPath = antisamename.GenerateUniquePath(newPath)

	var err error
	delay := config.RetryDelay
	for i := 0; i < config.MaxRetryAttempts; i++ {
		err = os.Rename(oldPath, newPath)
		if err == nil {
			return nil
		}
		if !IsFileBusyError(err) {
			break
		}
		time.Sleep(delay)
		delay *= 2
	}

	return fmt.Errorf("%s: %s → %s", i18n.Tr("rename_failed_format"), oldPath, newPath)

}
