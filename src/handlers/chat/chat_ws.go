package chat

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"sync"
	"time"
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
				panic(err)
			}
			msgJson, err := newMsgJson(WsConnect, uid.Value, chatId, map[string]string{})
			if err != nil {
				panic(err)
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
		err := json.Unmarshal(msgData, &msg)
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			return
		}
		db := conf.DatabaseI
		err = db.Connect()
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			return
		}
		actionFunc, ok := actionsMap[msg.Type]
		if ok {
			actionFunc(r, msg, db, &msgJson)
		}
		err = db.Close()
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			return
		}
		// Sending a message to a specific chat room.
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
	actionsAfterInsertNewMessage(r, msgJson, &messageData, &newMsg, db)
}

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
		panic(err)
	}
	newMsg["Images"] = messageData.Msg["Images"]

	actionsAfterInsertNewMessage(r, msgJson, &messageData, &newMsg, db)
}

func actionsAfterInsertNewMessage(r *http.Request, msgJson *string, messageData *Message, newMsg *map[string]string, db *database.Database) {
	// Increment msg count.
	err := IncrementChatMsgCountFromDb(r, messageData.ChatId, messageData.Uid, db)
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
		panic(err)
		//*msgJson = wsError(msg.Uid, msg.ChatId, decMsgCount.Error.Error())
		//return
	}
	*msgJson, err = newMsgJson(messageData.Type, messageData.Uid, messageData.ChatId, messageData.Msg)
	if err != nil {
		*msgJson = wsError(messageData.Uid, messageData.ChatId, err.Error())
		return
	}
}

// globalIncrementMessages If all conditions are met, sends a notification increase message to the notification socket.
func globalIncrementMessages(r *http.Request, sendUserId string, chatId string, db *database.Database) error {
	user, err := GetRecipientUser(chatId, sendUserId, db)
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
	var cm ChatMessage
	err = dbutils.FillStructFromDb(msg[0], &cm)
	if err != nil {
		return nil, err
	}
	msgMap := map[string]string{"Id": cm.Id, "UserId": cm.UserId, "Text": cm.Text, "Date": cm.Date, "IsRead": cm.IsRead}
	return msgMap, nil
}
