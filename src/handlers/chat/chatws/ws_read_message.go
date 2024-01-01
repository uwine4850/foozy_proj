package chatws

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"net/http"
)

// handleWsReadMsg processing a message read by a user.
// Changes in the database of message status to "read".
// Sending data about the message to the client.
func handleWsReadMsg(r *http.Request, messageData Message, db *database.Database, msgJson *string) {
	db.AsyncQ().AsyncUpdate("updMsg", "chat_msg", []dbutils.DbEquals{
		{
			Name:  "is_read",
			Value: true,
		},
	}, dbutils.WHEquals(map[string]interface{}{"id": messageData.Msg["Id"]}, "AND"))
	// Decrement message count
	db.AsyncQ().AsyncQuery("decMsgCount", "UPDATE `chat_msg_count` SET `count`= `count` - 1 WHERE user = ? AND chat = ? ;",
		messageData.Uid, messageData.ChatId)
	db.AsyncQ().Wait()
	updMsg, _ := db.AsyncQ().LoadAsyncRes("updMsg")
	if updMsg.Error != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, updMsg.Error.Error())
		return
	}
	decMsgCount, _ := db.AsyncQ().LoadAsyncRes("decMsgCount")
	if decMsgCount.Error != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, decMsgCount.Error.Error())
		return
	}
	err := globalDecrementMessages(r, messageData.Uid, messageData.ChatId, db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, decMsgCount.Error.Error())
		return
	}
	*msgJson, err = newMsgJson(messageData.Type, messageData.Uid, messageData.ChatId, messageData.Msg)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
}

// globalDecrementMessages sends a signal to the notification socket that the global number of messages should decrease.
func globalDecrementMessages(r *http.Request, readUID string, chatId string, db *database.Database) error {
	count, err := db.SyncQ().QB().Select("count", "chat_msg_count").
		Where("chat", "=", chatId, "AND", "user", "=", readUID, "AND count = 0").Ex()
	if err != nil {
		return err
	}
	if count != nil {
		err := notification.SendGlobalDecrementMsg(r, readUID)
		if err != nil {
			return err
		}
	}
	return nil
}
