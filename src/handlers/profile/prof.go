package profile

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/object"
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

type ProfView struct {
	object.ObjView
}

func (v *ProfView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	id, _ := manager.GetSlugParams("id")
	UID, _ := manager.GetUserContext("UID")
	db := conf.NewDb()
	err := db.Connect()
	if err != nil {
		router.ServerError(w, err.Error())
		return map[string]interface{}{}
	}
	isChatExist, err := chatExist(id, UID, db)
	if err != nil {
		router.ServerError(w, err.Error())
		return map[string]interface{}{}
	}
	return map[string]interface{}{"isChatExist": isChatExist}
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

func InitProfileView() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	view := object.ObjView{
		UserView:     &ProfView{},
		Name:         "user",
		TemplatePath: "src/templates/profile.html",
		DB:           conf.NewDb(),
		TableName:    "auth",
		FillStruct:   UserData{},
		Slug:         "id",
	}
	return view.Call
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
