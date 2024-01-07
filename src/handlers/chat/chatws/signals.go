package chatws

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy_proj/src/handlers/chat"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
)

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
	recipientUser, err := chat.GetRecipientUser(messageData.ChatId, messageData.Uid, db)
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

func SendDeleteMessage(r *http.Request, uid string, chatId string, msgId string) error {
	msgJson, err := newMsgJson(WsDeleteMessage, uid, chatId, map[string]string{"msgId": msgId})
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, msgJson, "ws://localhost:8000/chat-ws", true)
	if err != nil {
		return err
	}
	return nil
}

func SendUpdateMessage(r *http.Request, uid string, chatId string, msgData map[string]string) error {
	msgJson, err := newMsgJson(WsUpdateMessage, uid, chatId, msgData)
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, msgJson, "ws://localhost:8000/chat-ws", true)
	if err != nil {
		return err
	}
	return nil
}
