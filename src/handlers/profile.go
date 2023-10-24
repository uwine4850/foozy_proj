package handlers

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"log"
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
		panic(err)
	}

}

func RegisterPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		panic(err)
	}
	err = frm.ValidateCsrfToken()
	if err != nil {
		panic(err)
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"name", "username", "password", "confirm_pass"}, frm.GetApplicationForm())
	if !ok {
		manager.SetUserContext("error", "Some field not exist.")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	if fields["password"] != fields["confirm_pass"] {
		manager.SetUserContext("error", "The passwords don't match.")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	db := database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		log.Panicln(err.Error())
	}

	// Register new user.
	err = newAuth.RegisterUser(fields["username"], fields["password"])
	if err != nil {
		manager.SetUserContext("error", err.Error())
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		log.Panicln(err.Error())
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			log.Panicln(err.Error())
		}
		_, err = db.SyncQ().Update("auth", []dbutils.DbEquals{{"name", fields["name"]}},
			[]dbutils.DbEquals{{"id", id}})
		if err != nil {
			log.Panicln(err.Error())
		}
	} else {
		manager.SetUserContext("error", fmt.Sprintf("Username %s not exist.", fields["username"]))
		http.Redirect(w, r, "/register", http.StatusFound)
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
		panic(err)
	}
}

func SignInPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		panic(err)
	}
	err = frm.ValidateCsrfToken()
	if err != nil {
		manager.SetUserContext("error", err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"username", "password"}, frm.GetApplicationForm())
	if !ok {
		manager.SetUserContext("error", "Some field not exist.")
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	// Connect to database.
	db := database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// New auth.
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = newAuth.LoginUser(fields["username"], fields["password"])
	if err != nil {
		manager.SetUserContext("error", err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	// If the user exists set their ID in cookies.
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		manager.SetUserContext("error", err.Error())
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			log.Panicln(err)
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
		log.Panicln(err)
	}
	defer func(DatabaseI *database.Database) {
		err := DatabaseI.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(conf.DatabaseI)
	user, err := conf.DatabaseI.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", id}}, 1)
	if err != nil {
		log.Panicln(err)
	}

	// Render 404 if user not found
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("User not found"))
		if err != nil {
			log.Panicln(err)
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
