package profile

import (
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"strconv"
)

func SignIn(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	manager.SetTemplatePath("src/templates/auth/signin.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
	err := manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

type SignInForm struct {
	Username []string `form:"username"`
	Password []string `form:"password"`
}

func SignInPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm, err := utils.ParseForm(r)
	if err != nil {
		return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
	}
	var signInForm SignInForm
	err = form.FillStructFromForm(frm, form.NewFillableFormStruct(&signInForm), []string{})
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	// Connect to database.
	db := conf.NewDb()
	err = db.Connect()
	if err != nil {
		return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)

	// New auth.
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
	}
	err = newAuth.LoginUser(signInForm.Username[0], signInForm.Password[0])
	if err != nil {
		return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
	}
	// If the user exists set their ID in cookies.
	user, err := newAuth.UserExist(signInForm.Username[0])
	if err != nil {
		return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			return func() { router.RedirectError(w, r, "/sign-in", err.Error(), manager) }
		}
		cookie := &http.Cookie{
			Name:     "UID",
			Value:    strconv.Itoa(id),
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
	return func() {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}
