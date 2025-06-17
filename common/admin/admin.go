package admin

import (
	"golang.org/x/sys/windows"
)

// IsAdmin checks if the current process is running as administrator.
func IsAdmin() bool {
	sid, _ := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	token := windows.Token(0)
	isMember, _ := token.IsMember(sid)
	return isMember
}
