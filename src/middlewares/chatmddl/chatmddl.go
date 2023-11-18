package chatmddl

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"regexp"
	"strconv"
)

func ChatPermissionMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	re := regexp.MustCompile(`^/chat/\d+$`)
	if !re.MatchString(r.URL.Path) {
		return
	}
	uid, err := r.Cookie("UID")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return
		}
		panic(err)
	}
	chatId, _ := manager.GetSlugParams("id")
	db := conf.DatabaseI
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
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"id": chatId,
	}, "AND"), 1)
	if err != nil {
		panic(err)
	}
	if chat == nil {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	user1, err := dbutils.ParseInt(chat[0]["user1"])
	if err != nil {
		router.ServerError(w, err.Error())
		return
	}
	user2, err := dbutils.ParseInt(chat[0]["user2"])
	if err != nil {
		router.ServerError(w, err.Error())
		return
	}
	uidInt, _ := strconv.Atoi(uid.Value)
	if user1 != uidInt && user2 != uidInt {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
}
