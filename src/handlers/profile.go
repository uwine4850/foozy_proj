package handlers

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"strconv"
)

func Register(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	manager.SetTemplatePath("src/templates/auth/register.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
	err := manager.RenderTemplate(w, r)
	if err != nil {
		utils.ServerError(w, err.Error())
	}
}

func RegisterPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm, err := utils.ParseForm(r)
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"name", "username", "password", "confirm_pass"}, frm.GetApplicationForm())
	if !ok {
		utils.RedirectError(w, r, "/register", "Some field not exist.", manager)
		return
	}
	if fields["password"] != fields["confirm_pass"] {
		utils.RedirectError(w, r, "/register", "The passwords don't match.", manager)
		return
	}
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}

	// Register new user.
	err = newAuth.RegisterUser(fields["username"], fields["password"])
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			utils.RedirectError(w, r, "/register", err.Error(), manager)
			return
		}
		_, err = db.SyncQ().Update("auth", []dbutils.DbEquals{{"name", fields["name"]}},
			[]dbutils.DbEquals{{"id", id}})
		if err != nil {
			utils.RedirectError(w, r, "/register", err.Error(), manager)
			return
		}
	} else {
		utils.RedirectError(w, r, "/register", fmt.Sprintf("Username %s not exist.", fields["username"]), manager)
		return
	}
	http.Redirect(w, r, "/sign-in", http.StatusFound)
}

func SignIn(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	manager.SetTemplatePath("src/templates/auth/signin.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
	err := manager.RenderTemplate(w, r)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
}

func SignInPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm, err := utils.ParseForm(r)
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"username", "password"}, frm.GetApplicationForm())
	if !ok {
		utils.RedirectError(w, r, "/sign-in", "Some field not exist.", manager)
		return
	}
	// Connect to database.
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)

	// New auth.
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	err = newAuth.LoginUser(fields["username"], fields["password"])
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	// If the user exists set their ID in cookies.
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
			return
		}
		cookie := &http.Cookie{
			Name:     "UID",
			Value:    strconv.Itoa(id),
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func ProfileView(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	id, _ := manager.GetSlugParams("id")
	err := conf.DatabaseI.Connect()
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	defer func(DatabaseI *database.Database) {
		err := DatabaseI.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(conf.DatabaseI)
	user, err := conf.DatabaseI.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", id}}, 1)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}

	// Render 404 if user not found
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("User not found"))
		if err != nil {
			utils.ServerError(w, err.Error())
			return
		}
		return
	}
	userD := userData{
		Name:        dbutils.ParseString(user[0]["name"]),
		Username:    dbutils.ParseString(user[0]["username"]),
		Avatar:      dbutils.ParseString(user[0]["avatar"]),
		Description: dbutils.ParseString(user[0]["description"]),
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": userD})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
}

type userData struct {
	Name        string
	Username    string
	Avatar      string
	Description string
}
