package chat

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"strconv"
)

func LoadMessages(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	var error string
	chatId := r.URL.Query().Get("chatid")
	msgId := r.URL.Query().Get("msgid")
	_type := r.URL.Query().Get("msgtype")
	first := r.URL.Query().Get("first")
	handler := r.URL.Query().Get("handler")
	uid, err := r.Cookie("UID")
	if err != nil {
		panic(err)
		//error = uid.Value
	}

	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		return func() { error = err.Error() }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			panic(err)
			//sendJson(map[string]string{"err": error}, w)
		}
	}(db)

	// Check permission
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHOutput{
		QueryStr:  "id = ? AND user1 = ? OR user2 = ?",
		QueryArgs: []interface{}{chatId, uid.Value, uid.Value},
	}, 1)
	if err != nil {
		panic(err)
		//return func() { error = err.Error() }
	}
	if chat == nil {
		return func() { error = "Permission dined" }
	}

	// Load messages form database
	var messages []map[string]interface{}
	if first == "1" {
		_type = handler
	}
	switch _type {
	case "read":
		_messages, err := db.SyncQ().Query("SELECT * FROM `chat_msg` WHERE chat = ? AND id < ? "+
			"ORDER BY id DESC LIMIT "+strconv.Itoa(conf.LoadMessages), chatId, msgId)
		if err != nil {
			panic(err)
			//return func() { error = err.Error() }
		}
		messages = _messages
	case "notread":
		_messages, err := db.SyncQ().Query("SELECT * FROM `chat_msg` WHERE chat = ? AND id > ? "+
			" LIMIT "+strconv.Itoa(conf.LoadMessages), chatId, msgId)
		if err != nil {
			panic(err)
			//return func() { error = err.Error() }
		}
		messages = _messages
	}
	var chatMessages []ChatMessage
	for i := 0; i < len(messages); i++ {
		var m ChatMessage
		err := dbutils.FillStructFromDb(messages[i], &m)
		if err != nil {
			panic(err)
			//return func() { error = err.Error() }
		}
		chatMessages = append(chatMessages, m)
	}
	if error != "" {
		sendJson(map[string]string{"err": error}, w)
	} else {
		sendJson(map[string]interface{}{"messages": chatMessages, "chatId": chatId, "uid": uid.Value, "type": _type, "first": first}, w)
	}
	return func() {}
}

func sendJson(data interface{}, w http.ResponseWriter) {
	marshal, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		panic(err)
	}
}
