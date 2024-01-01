package profile

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
)

type UserData struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Username    string `db:"username"`
	Avatar      string `db:"avatar"`
	Description string `db:"description"`
}

func ProfileView(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	id, _ := manager.GetSlugParams("id")
	db := conf.NewDb()
	err := db.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(DatabaseI *database.Database) {
		err := DatabaseI.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)
	user, err := db.SyncQ().Select([]string{"*"}, "auth", dbutils.WHEquals(map[string]interface{}{"id": id}, "AND"), 1)
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
	manager.SetUserContext("subscribe_user_id", id)
	var fillUserData UserData
	err = dbutils.FillStructFromDb(user[0], &fillUserData)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	UID, _ := manager.GetUserContext("UID")
	isChatExist, err := chatExist(id, UID, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": fillUserData, "isChatExist": isChatExist})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func chatExist(id any, uid any, db *database.Database) (int, error) {
	res, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHOutput{
		QueryStr:  "user1 = ? AND user2 = ? OR user1 = ? AND user2 = ?",
		QueryArgs: []interface{}{id, uid, uid, id},
	}, 1)
	if err != nil {
		return -1, err
	}
	if res == nil {
		return -1, nil
	}
	parseInt, err := dbutils.ParseInt(res[0]["id"])
	if err != nil {
		return -1, err
	}
	return parseInt, nil
}

func GetUserDataById(id string, db *database.Database) (*UserData, error) {
	userData, err := db.SyncQ().QB().Select("*", "auth").Where("id", "=", id).Ex()
	if err != nil {
		return nil, err
	}
	if len(userData) == 0 {
		return nil, err
	}
	var user UserData
	err = dbutils.FillStructFromDb(userData[0], &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
