package admin

import (
	"sync"

	"golang.org/x/sys/windows"
)

var (
	isAdmin  bool
	initOnce sync.Once
)

func IsAdmin() bool {
	initOnce.Do(func() {
		sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
		if err != nil {
			return
		}

		token := windows.GetCurrentProcessToken()
		defer token.Close()

		isMember, err := token.IsMember(sid)
		if err != nil {
			return
		}

		isAdmin = isMember
	})
	return isAdmin
}
