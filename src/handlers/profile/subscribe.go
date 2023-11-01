package profile

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
)

func SubscribePost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	subscribeUserId, _ := manager.GetUserContext("subscribe_user_id")
	manager.DelUserContext("subscribe_user_id")
	UID, _ := manager.GetUserContext("UID")
	db := conf.DatabaseI
	err := db.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)

	res, err := db.SyncQ().Select([]string{"*"}, "subscribers", []dbutils.DbEquals{
		{"subscriber", UID},
		{"profile", subscribeUserId},
	}, 1)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	if res != nil {
		_, err := db.SyncQ().Delete("subscribers", []dbutils.DbEquals{{"id", res[0]["id"]}})
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
	} else {
		// Subscribing to a user
		_, err = db.SyncQ().Insert("subscribers", map[string]interface{}{"subscriber": UID, "profile": subscribeUserId})
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
	}
	return func() {
		http.Redirect(w, r, fmt.Sprintf("/prof/%v", subscribeUserId), http.StatusFound)
	}
}
