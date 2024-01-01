package chatws

import (
	"github.com/uwine4850/foozy/pkg/database"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// handleWsImageNsg image message processing.
// The message is processed almost in the same way as a text message, only here images are stored in the database.
func handleWsImageNsg(r *http.Request, messageData Message, db *database.Database, msgJson *string) {
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
	err = saveMessageImages(messageData.Msg["Images"], newMsg["Id"], db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
	newMsg["Images"] = messageData.Msg["Images"]

	err = SendPopUpMessageNotification(r, &messageData, db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}

	actionsAfterInsertNewMessage(r, msgJson, &messageData, &newMsg, db)
}

// saveMessageImages saving images from the message in the database.
func saveMessageImages(imagesPaths string, newMessageId string, db *database.Database) error {
	images := strings.Split(imagesPaths, "\\")
	queryKeys := make([]string, 0)
	for i := 0; i < len(images); i++ {
		db.AsyncQ().AsyncInsert("saveImg"+strconv.Itoa(i), "chat_msg_images", map[string]interface{}{
			"parent_msg": newMessageId,
			"path":       images[i],
		})
		queryKeys = append(queryKeys, "saveImg"+strconv.Itoa(i))
	}
	db.AsyncQ().Wait()
	for i := 0; i < len(queryKeys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(queryKeys[i])
		if ok {
			if res.Error != nil {
				return res.Error
			}
		}
	}
	return nil
}
