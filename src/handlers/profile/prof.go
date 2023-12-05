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
	db := conf.DatabaseI
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
	isSubscribe, err := userIsSubscribe(id, UID, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	subCount, err := getCountSubscribers(id, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	isChatExist, err := chatExist(id, UID, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": fillUserData, "isSubscribe": isSubscribe, "subCount": subCount, "isChatExist": isChatExist})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func userIsSubscribe(subscribeUserId any, uid any, db *database.Database) (bool, error) {
	res, err := db.SyncQ().Select([]string{"*"}, "subscribers", dbutils.WHEquals(map[string]interface{}{
		"subscriber": uid,
		"profile":    subscribeUserId,
	}, "AND"), 1)
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func getCountSubscribers(profileId any, db *database.Database) (int, error) {
	count, err := db.SyncQ().Count([]string{"*"}, "subscribers", dbutils.WHEquals(map[string]interface{}{
		"profile": profileId,
	}, "AND"), 0)
	if err != nil {
		return 0, err
	}
	parseInt, err := dbutils.ParseInt(count[0]["COUNT(*)"])
	if err != nil {
		return 0, err
	}
	return parseInt, nil
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
