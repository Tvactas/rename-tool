package filestatus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	baseNewPath := newPath
	counter := 1
	ext := filepath.Ext(baseNewPath)
	name := baseNewPath[:len(baseNewPath)-len(ext)]
	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		newPath = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}

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
