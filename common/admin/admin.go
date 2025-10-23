package admin

import (
	"rename-tool/common/applog"
	"rename-tool/setting/i18n"

	"golang.org/x/sys/windows"
)

// Check IsAdmin
func IsAdmin() bool {
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		applog.Logger.Printf("[ADMIN ERROR]  %s: %v", i18n.LogTr("CreateWellKnownSidFail"), err)
		return false
	}

	token := windows.GetCurrentProcessToken()

	isMember, err := token.IsMember(sid)
	if err != nil {
		applog.Logger.Printf("[ADMIN ERROR]  %s: %v", i18n.LogTr("CheckIsMember"), err)
		return false
	}

	applog.Logger.Printf("[ADMIN LOGIN]  %s: %v", i18n.LogTr("LoginIdentity"), isMember)
	return isMember
}
