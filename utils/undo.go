package utils

import (
	"fmt"
	"os"

	"rename-tool/setting/global"
)

// UndoRename handles undoing previous rename operations in memory
func UndoRename() {
	if len(global.Logs) == 0 {
		warningDiaLog(global.MainWindow, dialogTr("noUndoOperations"))
		return
	}

	var (
		newLogs      []global.RenameLog // 保留未撤销的日志
		busyFiles    []string           // 无法撤销的文件
		successCount int
	)

	// 倒序遍历日志，最新的重命名先撤销
	for i := len(global.Logs) - 1; i >= 0; i-- {
		log := global.Logs[i]

		// 判断目标文件是否存在（即要撤销的“新文件名”）
		if _, err := os.Stat(log.New); err == nil {
			// 尝试把文件名改回原名
			if err := os.Rename(log.New, log.Original); err == nil {
				successCount++
				continue // 撤销成功，不保留这条日志
			} else {
				// 文件被占用或权限问题
				busyFiles = append(busyFiles, log.New)
				newLogs = append([]global.RenameLog{log}, newLogs...)
			}
		} else {
			// 新文件不存在，说明用户手动删了或改了名
			newLogs = append([]global.RenameLog{log}, newLogs...)
		}
	}

	// 更新全局日志（只保留未撤销成功的）
	global.Logs = newLogs

	// 反馈结果
	switch {
	case successCount == 0 && len(busyFiles) == 0:
		warningDiaLog(global.MainWindow, dialogTr("noUndoOperations"))

	case len(busyFiles) > 0:
		warningMultiDiaLog(global.MainWindow, busyFiles)

	default:
		successDiaLog(global.MainWindow, fmt.Sprintf(dialogTr("undoSuccess"), successCount))
	}
}
