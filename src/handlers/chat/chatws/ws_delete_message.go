package chatws

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"net/http"
)

func handleWsDeleteMessage(r *http.Request, messageData Message, db *database.Database, msgJson *string) {
	msgId := messageData.Msg["msgId"]
	deleteEquals := map[string]interface{}{
		"id":   msgId,
		"user": messageData.Uid,
		"chat": messageData.ChatId,
	}
	// Check permissions
	permissionsOk, err := db.SyncQ().Select([]string{"*"}, "chat_msg", dbutils.WHEquals(deleteEquals, "AND"), 1)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
	if permissionsOk == nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, "the user cannot delete this message")
		return
	}
	// Delete message
	_, err = db.SyncQ().Delete("chat_msg", dbutils.WHEquals(deleteEquals, "AND"))
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, "the user cannot delete this message")
		return
	}
	*msgJson, err = newMsgJson(WsDeleteMessage, messageData.Uid, messageData.ChatId, messageData.Msg)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
}
