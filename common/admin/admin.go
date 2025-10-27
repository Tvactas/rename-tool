package admin

import (
	"golang.org/x/sys/windows"
)

// Check IsAdmin
func IsAdmin() bool {
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		logEvent("ADMIN ERROR", "CreateWellKnownSidFail", err)
		return false
	}

	token := windows.GetCurrentProcessToken()

	isMember, err := token.IsMember(sid)
	if err != nil {
		logEvent("ADMIN ERROR", "FailCheckMember", err)

		return false
	}
	logEvent("ADMIN LOGIN", "LoginIdentity", isMember)

	return isMember
}
