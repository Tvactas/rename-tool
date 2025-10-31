package filestatus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rename-tool/common/antisamename"
	"rename-tool/setting/config"
	"rename-tool/setting/i18n"
	"strings"
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

// SafeCaseOnlyRename performs a safe two-phase rename when the only difference
// between oldPath and newPath is letter case (Windows case-insensitive).
// It renames oldPath -> tempUnique -> newPath to ensure the case change takes effect.
func SafeCaseOnlyRename(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}

	// Only perform two-phase if paths are equal ignoring case
	if !strings.EqualFold(oldPath, newPath) {
		// Not a case-only change; fall back to normal rename with retry/avoid overwrite logic
		return RenameFile(oldPath, newPath)
	}

	dir := filepath.Dir(oldPath)
	base := filepath.Base(oldPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	// Construct a temp unique path in the same directory
	tempCandidate := filepath.Join(dir, fmt.Sprintf("%s.tmp-%d%s", name, time.Now().UnixNano(), ext))
	tempPath := antisamename.GenerateUniquePath(tempCandidate)

	// Phase 1: old -> temp (with retry on file-busy)
	if err := renameWithRetry(oldPath, tempPath); err != nil {
		return fmt.Errorf("%s: %s → %s", i18n.Tr("rename_failed_format"), oldPath, tempPath)
	}

	// Phase 2: temp -> new (with retry on file-busy)
	if err := renameWithRetry(tempPath, newPath); err != nil {
		// best-effort: try to move back to original name
		_ = renameWithRetry(tempPath, oldPath)
		return fmt.Errorf("%s: %s → %s", i18n.Tr("rename_failed_format"), oldPath, newPath)
	}

	return nil
}

// renameWithRetry attempts os.Rename with exponential backoff on file-busy errors.
func renameWithRetry(from, to string) error {
	delay := config.RetryDelay
	for i := 0; i < config.MaxRetryAttempts; i++ {
		if err := os.Rename(from, to); err != nil {
			if IsFileBusyError(err) {
				time.Sleep(delay)
				delay *= 2
				continue
			}
			return err
		}
		return nil
	}
	return fmt.Errorf("rename retries exhausted: %s → %s", from, to)
}
