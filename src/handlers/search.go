package handlers

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	manager.SetTemplatePath("src/templates/search.html")
	err := manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

func SearchHandlerPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { sendError(w, err) }
	}
	searchUsername := newForm.Value("search")
	db := conf.NewDb()
	err = db.Connect()
	if err != nil {
		return func() { sendError(w, err) }
	}
	usersDb, err := db.SyncQ().QB().Select("*", "auth").
		Where("username LIKE \"%" + searchUsername + "%\" LIMIT 5").Ex()
	if err != nil {
		return func() { sendError(w, err) }
	}
	// Close database
	err = db.Close()
	if err != nil {
		return func() { sendError(w, err) }
	}
	var users []profile.UserData
	for i := 0; i < len(usersDb); i++ {
		var user profile.UserData
		err := dbutils.FillStructFromDb(usersDb[i], &user)
		if err != nil {
			return func() { sendError(w, err) }
		}
		users = append(users, user)
	}
	jsonUsers, err := json.Marshal(map[string]interface{}{"users": users})
	if err != nil {
		return func() { sendError(w, err) }
	}
	w.Write(jsonUsers)
	return func() {}
}

func sendError(w http.ResponseWriter, err error) {
	jsonUsers, err := json.Marshal(map[string]interface{}{"error": err.Error()})
	if err != nil {
		panic(err)
	}
	w.Write(jsonUsers)
}
