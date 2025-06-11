package admin

import (
	"os"
	"os/exec"
	"syscall"
)

// IsAdmin 检查当前是否以管理员权限运行
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// RunAsAdmin 以管理员权限重新运行程序
func RunAsAdmin() {
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := os.Args[1:]

	// 在Windows上请求管理员权限
	cmd := exec.Command("cmd", "/C", "start", "runas", exe)
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = cwd
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start() // 忽略错误，用户可能取消UAC提示
}
