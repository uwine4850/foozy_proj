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
	manager.SetUserContext("subscribe_user_id", id)
	var fillUserData userData
	err = dbutils.FillStructFromDb(user[0], &fillUserData)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	UID, _ := manager.GetUserContext("UID")
	isSubscribe, err := userIsSubscribe(id, UID, conf.DatabaseI)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	subCount, err := getCountSubscribers(id, conf.DatabaseI)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": fillUserData, "isSubscribe": isSubscribe, "subCount": subCount})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func userIsSubscribe(subscribeUserId any, uid any, db interfaces.IDatabase) (bool, error) {
	res, err := db.SyncQ().Select([]string{"*"}, "subscribers", []dbutils.DbEquals{
		{"subscriber", uid},
		{"profile", subscribeUserId},
	}, 1)
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func getCountSubscribers(profileId any, db interfaces.IDatabase) (int, error) {
	res, err := db.SyncQ().Query("SELECT COUNT(*) FROM `subscribers` WHERE profile = ?;", profileId)
	if err != nil {
		return 0, err
	}
	parseInt, err := dbutils.ParseInt(res[0]["COUNT(*)"])
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}
