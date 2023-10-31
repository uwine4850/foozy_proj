package profile

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
)

func ProfileLogOutPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	cookie := &http.Cookie{
		Name:     "UID",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
	return func() {
		http.Redirect(w, r, "/sign-in", http.StatusFound)
	}
}
