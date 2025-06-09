package FileStatus

import "strings"

// ============== 文件占用处理函数 ==============
func IsFileBusyError(err error) bool {
	if err == nil {
		return false
	}

	// 检查常见的文件占用错误
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "access is denied") ||
		strings.Contains(errMsg, "file is locked") ||
		strings.Contains(errMsg, "process cannot access the file") ||
		strings.Contains(errMsg, "the file is being used by another process")
}
