package notificationmddl

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
)

func NotificationCountMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManagerData) {
	if utils.SliceContains([]string{"/notification-ws", "/chat-ws", "/load-messages"}, r.URL.Path) {
		return
	}
	uid, err := r.Cookie("UID")
	if err != nil {
		return
	}
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	count, err := db.SyncQ().QB().Select("count", "chat_msg_count").
		Where("user", "=", uid.Value, "AND", "count > 0").Ex()
	if err != nil {
		return
	}
	manager.SetContext(map[string]interface{}{"msgCount": len(count)})
}
