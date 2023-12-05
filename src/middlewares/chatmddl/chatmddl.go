package chatmddl

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/middlewares"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"regexp"
	"strconv"
)

func ChatPermissionMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData) {
	re := regexp.MustCompile(`^/chat/\d+$`)
	if !re.MatchString(r.URL.Path) {
		return
	}
	uid, err := r.Cookie("UID")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			middlewares.SkipNextPageAndRedirect(manager, w, r, "/home")
			return
		}
		middlewares.SetMddlError(err, manager)
		return
	}
	chatId, _ := manager.GetSlugParams("id")
	db := conf.DatabaseI
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
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"id": chatId,
	}, "AND"), 1)
	if err != nil {
		middlewares.SetMddlError(err, manager)
		return
	}
	if chat == nil {
		router.ServerForbidden(w)
		middlewares.SkipNextPage(manager)
		return
	}
	user1, err := dbutils.ParseInt(chat[0]["user1"])
	if err != nil {
		middlewares.SetMddlError(err, manager)
		return
	}
	user2, err := dbutils.ParseInt(chat[0]["user2"])
	if err != nil {
		middlewares.SetMddlError(err, manager)
		return
	}
	uidInt, _ := strconv.Atoi(uid.Value)
	if user1 != uidInt && user2 != uidInt {
		router.ServerForbidden(w)
		middlewares.SkipNextPage(manager)
		return
	}
}
