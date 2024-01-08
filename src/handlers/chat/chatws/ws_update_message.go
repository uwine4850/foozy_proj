package chatws

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy_proj/src/handlers/chat"
	"net/http"
	"strings"
)

func handleWsUpdateMessage(r *http.Request, messageData Message, db *database.Database, msgJson *string) {
	// Update text
	db.AsyncQ().QB("updText").Update("chat_msg", map[string]interface{}{
		"text": messageData.Msg["text"],
	}).Where("id", "=", messageData.Msg["id"]).Ex()

	// Delete images
	images := strings.Split(messageData.Msg["delImages"], "\\")
	var delImageId []interface{}
	for i := 0; i < len(images); i++ {
		if images[i] == "" {
			continue
		}
		imgId, err := db.SyncQ().QB().Select("id", "chat_msg_images").Where("path", "=", images[i]).Ex()
		if err != nil {
			*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
			return
		}
		delImageId = append(delImageId, imgId[0]["id"])
	}
	for i := 0; i < len(delImageId); i++ {
		_, err := db.SyncQ().Delete("chat_msg_images", dbutils.WHEquals(map[string]interface{}{
			"id": delImageId[i],
		}, "AND"))
		if err != nil {
			*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
			return
		}
	}

	db.AsyncQ().Wait()
	updText, _ := db.AsyncQ().LoadAsyncRes("updText")
	if updText.Error != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, updText.Error.Error())
		return
	}

	// Retrieving an updated message. If it does not have images and text, delete this message.
	chatMessage, err := db.SyncQ().QB().Select("*", "chat_msg").Where("id", "=", messageData.Msg["id"]).Ex()
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, updText.Error.Error())
		return
	}
	messageImages, err := chat.LoadMessageImages(messageData.ChatId, db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, updText.Error.Error())
		return
	}
	var fillMessage chat.ChatMessage
	err = dbutils.FillStructFromDb(chatMessage[0], &fillMessage)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, updText.Error.Error())
		return
	}
	fillMessage.Images = messageImages

	if fillMessage.Text == "" && len(fillMessage.Images) == 0 {
		err := SendDeleteMessage(r, messageData.Uid, messageData.ChatId, fillMessage.Id)
		if err != nil {
			*msgJson = wsError(messageData.Uid, messageData.ChatId, updText.Error.Error())
			return
		}
		return
	}

	// If the message is not deleted, a signal is sent that it has been updated.
	*msgJson, err = newMsgJson(WsUpdateMessage, messageData.Uid, messageData.ChatId, messageData.Msg)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}

}
