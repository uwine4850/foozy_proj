package profilemddl

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/middlewares"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
)

func AuthMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData) {
	if r.Header.Get("Connection") == "Upgrade" {
		return
	}
	if r.URL.Path != "/register" && r.URL.Path != "/sign-in" && r.URL.Path != "/sign-in-post" && r.URL.Path != "/register-post" {
		uid, err := r.Cookie("UID")
		if err != nil {
			middlewares.SkipNextPageAndRedirect(manager, w, r, "/sign-in")
			return
		}
		// Connect to database.
		db := database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")
		err = db.Connect()
		if err != nil {
			middlewares.SetMddlError(err, manager)
			return
		}
		defer func(db *database.Database) {
			err := db.Close()
			if err != nil {
				middlewares.SetMddlError(err, manager)
				return
			}
		}(db)

		user, err := db.SyncQ().Select([]string{"*"}, "auth", dbutils.WHEquals(map[string]interface{}{"id": uid.Value}, "AND"), 1)
		if err != nil {
			middlewares.SetMddlError(err, manager)
			return
		}
		if user == nil {
			middlewares.SkipNextPageAndRedirect(manager, w, r, "/sign-in")
			return
		}
		var currentUser profile.UserData
		err = dbutils.FillStructFromDb(user[0], &currentUser)
		if err != nil {
			middlewares.SetMddlError(err, manager)
			return
		}
		manager.SetContext(map[string]interface{}{"UID": uid.Value})
		manager.SetContext(map[string]interface{}{"currentUser": currentUser})
		manager.SetUserContext("UID", uid.Value)
		manager.SetUserContext("currentUser", currentUser)
	}
}
