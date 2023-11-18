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

func MySubscriptions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	uid, ok := manager.GetUserContext("UID")
	if !ok {
		return func() { router.ServerError(w, "err") }
	}
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
	mySubscriptionsId, err := db.SyncQ().Select([]string{"profile"}, "subscribers", dbutils.WHEquals(map[string]interface{}{
		"subscriber": uid,
	}, "AND"), 0)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	var idList []interface{}
	for i := 0; i < len(mySubscriptionsId); i++ {
		idList = append(idList, fmt.Sprintf("%v", mySubscriptionsId[i]["profile"]))
	}
	usersData := make([]UserData, 0)
	if idList != nil {
		fmt.Println(idList)
		users, err := db.SyncQ().Select([]string{"*"}, "auth", dbutils.WHInSlice(map[string][]interface{}{
			"id": idList,
		}, "AND"), 0)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		for i := 0; i < len(users); i++ {
			fill := UserData{}
			err := dbutils.FillStructFromDb(users[i], &fill)
			if err != nil {
				return func() { router.ServerError(w, err.Error()) }
			}
			usersData = append(usersData, fill)
		}
	}
	manager.SetContext(map[string]interface{}{"subscriptions": usersData})
	manager.SetTemplatePath("src/templates/my_subscriptions.html")
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}
