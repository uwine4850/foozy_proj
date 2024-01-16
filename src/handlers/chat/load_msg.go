package chat

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"strconv"
)

func LoadMessages(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	chatId := r.URL.Query().Get("chatid")
	msgId := r.URL.Query().Get("msgid")
	_type := r.URL.Query().Get("msgtype")
	first := r.URL.Query().Get("first")
	handler := r.URL.Query().Get("handler")
	uid, err := r.Cookie("UID")
	if err != nil {
		return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
	}

	db := conf.NewDb()
	err = db.Connect()
	if err != nil {
		return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.SendJson(map[string]string{"err": err.Error()}, w)
		}
	}(db)

	// Check permission
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHOutput{
		QueryStr:  "id = ? AND user1 = ? OR user2 = ?",
		QueryArgs: []interface{}{chatId, uid.Value, uid.Value},
	}, 1)
	if err != nil {
		return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
	}
	if chat == nil {
		return func() { router.SendJson(map[string]string{"err": "Permission dined"}, w) }
	}

	// If this is the first message the message type will be equal to the message type.
	// This is done to start loading both read and unread messages.
	if first == "1" {
		_type = handler
	}
	// Load messages form database
	var messages []map[string]interface{}
	switch _type {
	case "read":
		_messages, err := db.SyncQ().Query("SELECT * FROM `chat_msg` WHERE chat = ? AND id < ? "+
			"ORDER BY id DESC LIMIT "+strconv.Itoa(conf.LoadMessages), chatId, msgId)
		if err != nil {
			return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
		}
		messages = _messages
	case "notread":
		_messages, err := db.SyncQ().Query("SELECT * FROM `chat_msg` WHERE chat = ? AND id > ? "+
			" LIMIT "+strconv.Itoa(conf.LoadMessages), chatId, msgId)
		if err != nil {
			return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
		}
		messages = _messages
	}
	// Filling the []ChatMessage slice with message data.
	var chatMessages []ChatMessage
	for i := 0; i < len(messages); i++ {
		var m ChatMessage
		err := dbutils.FillStructFromDb(messages[i], &m)
		if err != nil {
			return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
		}
		images, err := LoadMessageImages(m.Id, db)
		if err != nil {
			return func() { router.SendJson(map[string]string{"err": err.Error()}, w) }
		}
		m.Images = images
		chatMessages = append(chatMessages, m)
	}
	_ = router.SendJson(map[string]interface{}{"messages": chatMessages, "chatId": chatId, "uid": uid.Value, "type": _type, "first": first}, w)
	return func() {}
}
