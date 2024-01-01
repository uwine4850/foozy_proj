package chatws

import (
	"github.com/uwine4850/foozy/pkg/database"
	"net/http"
	"time"
)

// handleWsTextMsg processing a message sent to the chat room.
// The message is saved to the database, then parsed and sent back to the client.
func handleWsTextMsg(r *http.Request, messageData Message, db *database.Database, msgJson *string) {
	if messageData.Msg["Text"] == "" {
		return
	}
	newMsgData := map[string]interface{}{
		"user":    messageData.Uid,
		"chat":    messageData.ChatId,
		"text":    messageData.Msg["Text"],
		"date":    time.Now(),
		"is_read": false,
	}
	_, err := db.SyncQ().Insert("chat_msg", newMsgData)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
	newMsg, err := getNewMsg(db, newMsgData)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}

	// Send popup message
	err = SendPopUpMessageNotification(r, &messageData, db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}

	actionsAfterInsertNewMessage(r, msgJson, &messageData, &newMsg, db)
}
