package admin

import (
	"golang.org/x/sys/windows"
)

// 判断是否为管理员权限打开
// 暴露IsAdmin函数
// 在主页mainWindows中显示
func IsAdmin() bool {
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		logEvent("ADMIN ERROR", "createWellKnownSidFail", err)
		return false
	}

	token := windows.GetCurrentProcessToken()

	isMember, err := token.IsMember(sid)
	if err != nil {
		logEvent("ADMIN ERROR", "failCheckMember", err)

		return false
	}
	logEvent("ADMIN LOGIN", "loginIdentity", isMember)

	return isMember
}
