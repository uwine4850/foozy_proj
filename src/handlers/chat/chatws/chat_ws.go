package chatws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/chat"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"sync"
)

const (
	WsConnect = iota
	WsTextMsg
	WsReadMsg
	WsError
	WsImageNsg
)

type Message struct {
	Type   int
	Uid    string
	ChatId string
	Msg    map[string]string
}

type ActionFunc func(r *http.Request, messageData Message, db *database.Database, msgJson *string)

var actionsMap = map[int]ActionFunc{
	WsTextMsg:  handleWsTextMsg,
	WsReadMsg:  handleWsReadMsg,
	WsImageNsg: handleWsImageNsg,
}

var chatConnections = make(map[string][]*websocket.Conn)
var connections = make(map[*websocket.Conn]string)

func WsHandler(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	var mu sync.Mutex
	chatId, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, "Id chat was not found.") }
	}
	ws := manager.CurrentWebsocket()
	ws.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		once := notification.GetRequestOnce(r)
		if once == "false" {
			connChatId := connections[conn]
			chatConnections[connChatId] = utils.RemoveElement(chatConnections[connChatId], conn)
			delete(connections, conn)
		}
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	})
	ws.OnConnect(onConnect(r, ws, chatId))
	ws.OnMessage(onMessage(r, ws, &mu))
	err := ws.ReceiveMessages(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

// onConnect the function is executed when a new user joins the chat room.
// The variable chatConnections records information about the chat and the users who are members of it, for example, chatId = conn.
// The connections slice contains information about each connection and its chat, for example, conn = chatId.
// After all this, a message is sent to the client.
func onConnect(r *http.Request, ws interfaces.IWebsocket, chatId string) func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	return func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		once := notification.GetRequestOnce(r)
		if once == "false" {
			connections[conn] = chatId
			chatConnections[chatId] = append(chatConnections[chatId], conn)
			uid, err := r.Cookie("UID")
			if err != nil {
				ws.SendMessage(websocket.TextMessage, []byte(wsError(uid.Value, chatId, err.Error())), conn)
			}
			msgJson, err := newMsgJson(WsConnect, uid.Value, chatId, map[string]string{})
			if err != nil {
				ws.SendMessage(websocket.TextMessage, []byte(wsError(uid.Value, chatId, err.Error())), conn)
			}
			err = ws.SendMessage(websocket.TextMessage, []byte(msgJson), conn)
			if err != nil {
				panic(err)
			}
		}
	}
}

// onMessage processes messages from the user.
// Determines the type of message and then calls the appropriate handler.
// After processing, sends the message data back to the client.
func onMessage(r *http.Request, ws interfaces.IWebsocket, mu *sync.Mutex) func(messageType int, msgData []byte, conn *websocket.Conn) {
	return func(messageType int, msgData []byte, conn *websocket.Conn) {
		mu.Lock()
		defer mu.Unlock()
		var msg Message
		var msgJson string
		var db *database.Database
		var actionFunc ActionFunc
		var ok bool

		err := json.Unmarshal(msgData, &msg)
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			goto sendMessage
		}
		db = conf.NewDb()
		err = db.Connect()
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			goto sendMessage
		}
		actionFunc, ok = actionsMap[msg.Type]
		if ok {
			actionFunc(r, msg, db, &msgJson)
		}
		err = db.Close()
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		}
		// Sending a message to a specific chat room.
	sendMessage:
		for i := 0; i < len(chatConnections[msg.ChatId]); i++ {
			if msgJson == "" {
				return
			}
			err = ws.SendMessage(websocket.TextMessage, []byte(msgJson), chatConnections[msg.ChatId][i])
			if err != nil {
				panic(err)
			}
		}
	}
}

// actionsAfterInsertNewMessage function is used to start similar actions after sending a message.
func actionsAfterInsertNewMessage(r *http.Request, msgJson *string, messageData *Message, newMsg *map[string]string, db *database.Database) {
	// Increment msg count.
	err := chat.IncrementChatMsgCountFromDb(r, messageData.ChatId, messageData.Uid, messageData.Msg["Text"], db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
	err = globalIncrementMessages(r, messageData.Uid, messageData.ChatId, db)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
	*msgJson, err = newMsgJson(messageData.Type, messageData.Uid, messageData.ChatId, *newMsg)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
}

// globalIncrementMessages If all conditions are met, sends a notification increase message to the notification socket.
func globalIncrementMessages(r *http.Request, sendUserId string, chatId string, db *database.Database) error {
	user, err := chat.GetRecipientUser(chatId, sendUserId, db)
	if err != nil {
		return err
	}
	count, err := db.SyncQ().QB().Select("count", "chat_msg_count").
		Where("chat", "=", chatId, "AND", "user", "=", user.Id, "AND count = 1").Ex()
	if err != nil {
		return err
	}
	if count != nil {
		err := notification.SendGlobalIncrementMsg(r, user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// getNewMsg retrieves the last saved message in the database.
func getNewMsg(db *database.Database, msgData map[string]interface{}) (map[string]string, error) {
	delete(msgData, "date")
	equals := dbutils.WHEquals(msgData, "AND")
	equals.QueryStr += " ORDER BY Id DESC "
	msg, err := db.SyncQ().Select([]string{"*"}, "chat_msg", equals, 1)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("no new message found")
	}
	var cm chat.ChatMessage
	err = dbutils.FillStructFromDb(msg[0], &cm)
	if err != nil {
		return nil, err
	}
	msgMap := map[string]string{"Id": cm.Id, "UserId": cm.UserId, "Text": cm.Text, "Date": cm.Date, "IsRead": cm.IsRead}
	return msgMap, nil
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
