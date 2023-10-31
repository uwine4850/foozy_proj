package profile

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
)

type userData struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Username    string `db:"username"`
	Avatar      string `db:"avatar"`
	Description string `db:"description"`
}

func ProfileView(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	id, _ := manager.GetSlugParams("id")
	err := conf.DatabaseI.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(DatabaseI *database.Database) {
		err := DatabaseI.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(conf.DatabaseI)
	user, err := conf.DatabaseI.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", id}}, 1)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	// Render 404 if user not found
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("User not found"))
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		return func() {}
	}
	var fillUserData userData
	err = dbutils.FillStructFromDb(user[0], &fillUserData)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": fillUserData})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}
