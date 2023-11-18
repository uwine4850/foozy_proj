package profile

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	manager.SetTemplatePath("src/templates/auth/register.html")
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

type RegisterForm struct {
	Name        []string `form:"name"`
	Username    []string `form:"username"`
	Password    []string `form:"password"`
	ConfirmPass []string `form:"confirm_pass"`
}

func RegisterPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm, err := utils.ParseForm(r)
	if err != nil {
		return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
	}
	var registerForm RegisterForm
	err = form.FillStructFromForm(frm, &registerForm, []string{})
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	if registerForm.Password[0] != registerForm.ConfirmPass[0] {
		return func() { router.RedirectError(w, r, "/register", "The passwords don't match.", manager) }
	}
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
	}

	// Register new user.
	err = newAuth.RegisterUser(registerForm.Username[0], registerForm.Password[0])
	if err != nil {
		return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
	}
	user, err := newAuth.UserExist(registerForm.Username[0])
	if err != nil {
		return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
		}
		_, err = db.SyncQ().Update("auth", []dbutils.DbEquals{{"name", registerForm.Name[0]}},
			dbutils.WHEquals(map[string]interface{}{"id": id}, "AND"))
		if err != nil {
			return func() { router.RedirectError(w, r, "/register", err.Error(), manager) }
		}
	} else {
		return func() {
			router.RedirectError(w, r, "/register", fmt.Sprintf("Username %s not exist.", registerForm.Username[0]), manager)
		}
	}
	return func() {
		http.Redirect(w, r, "/sign-in", http.StatusFound)
	}
}
