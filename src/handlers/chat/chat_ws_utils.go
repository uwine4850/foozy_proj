package chat

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"strconv"
	"strings"
)

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

// wsError sending the error to the client.
func wsError(uid string, chatId string, error string) string {
	msgJson, err := newMsgJson(WsError, uid, chatId, map[string]string{"Error": error})
	if err != nil {
		panic(err)
	}
	return msgJson
}

// newMsgJson sending a response to the client in json format.
func newMsgJson(_type int, uid string, chatId string, msg map[string]string) (string, error) {
	m := Message{
		Type:   _type,
		Uid:    uid,
		ChatId: chatId,
		Msg:    msg,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// IncrementChatMsgCountFromDb Increases the number of unread messages of a specific user in a specific chat by 1.
// sendUid - id of the user who sent the message.
func IncrementChatMsgCountFromDb(r *http.Request, chatId string, sendUid string, db *database.Database) error {
	recipientUser, err := GetRecipientUser(chatId, sendUid, db)
	if err != nil {
		return err
	}
	_, err = db.SyncQ().Query("UPDATE `chat_msg_count` SET `count`= `count` + 1 WHERE user = ? AND chat = ? ;", recipientUser.Id, chatId)
	if err != nil {
		return err
	}
	err = notification.SendIncrementMsgChatCount(r, recipientUser.Id, chatId)
	if err != nil {
		return err
	}
	return nil
}

// SendTextMessage sends a new text message alert to the chat websocket.
func SendTextMessage(r *http.Request, msg *Message) error {
	msgJson, err := newMsgJson(WsTextMsg, msg.Uid, msg.ChatId, msg.Msg)
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, msgJson, "ws://localhost:8000/chat-ws", true)
	if err != nil {
		return err
	}
	return nil
}

// SendImageMessage sends a new image message alert to the chat websocket.
func SendImageMessage(r *http.Request, msg *Message) error {
	msgJson, err := newMsgJson(WsImageNsg, msg.Uid, msg.ChatId, msg.Msg)
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, msgJson, "ws://localhost:8000/chat-ws", true)
	if err != nil {
		return err
	}
	return nil
}

func SendPopUpMessageNotification(r *http.Request, messageData *Message, db *database.Database) error {
	user, err := profile.GetUserDataById(messageData.Uid, db)
	if err != nil {
		return err
	}
	recipientUser, err := GetRecipientUser(messageData.ChatId, messageData.Uid, db)
	data := make(map[string]string, len(messageData.Msg))
	for key, value := range messageData.Msg {
		data[key] = value
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	data["User"] = string(userBytes)
	err = notification.SendPopUpMessage(r, recipientUser.Id, &data)
	if err != nil {
		return err
	}
	return nil
}
