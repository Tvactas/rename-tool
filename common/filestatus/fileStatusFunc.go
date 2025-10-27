package filestatus

import "strings"

// isKnownFileBusyMessage checks if the message contains known "file busy" patterns.
func isKnownFileBusyMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "process cannot access the file") ||
		strings.Contains(msg, "file is being used by another process") ||
		strings.Contains(msg, "sharing violation") ||
		strings.Contains(msg, "file is locked")
}
