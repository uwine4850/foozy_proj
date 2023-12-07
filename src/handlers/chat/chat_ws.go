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
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	TypeConnect = iota
	TypeTextMsg
	TypeReadMsg
	TypeError
)

type Msg struct {
	Type   int
	Uid    string
	ChatId string
	Msg    map[string]string
}

var chatConnections = make(map[string][]*websocket.Conn)
var connections = make(map[*websocket.Conn]string)

func ChatWs(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	var mu sync.Mutex
	chatId, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, "Id chat was not found.") }
	}

	ws := manager.GetWebSocket()
	ws.OnClientClose(func(conn *websocket.Conn) {
		connChatId := connections[conn]
		chatConnections[connChatId] = utils.RemoveElement(chatConnections[connChatId], conn)
		delete(connections, conn)
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	})
	ws.OnConnect(onConnect(r, ws, chatId))
	ws.OnMessage(onMessage(ws, &mu))
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
func onConnect(r *http.Request, ws interfaces.IWebsocket, chatId string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		connections[conn] = chatId
		chatConnections[chatId] = append(chatConnections[chatId], conn)
		uid, err := r.Cookie("UID")
		if err != nil {
			panic(err)
		}
		msgJson, err := newMsgJson(TypeConnect, uid.Value, chatId, map[string]string{})
		if err != nil {
			panic(err)
		}
		err = ws.SendMessage(websocket.TextMessage, []byte(msgJson), conn)
		if err != nil {
			panic(err)
		}
	}
}

// onMessage processes messages from the user.
// Determines the type of message and then calls the appropriate handler.
// After processing, sends the message data back to the client.
func onMessage(ws interfaces.IWebsocket, mu *sync.Mutex) func(messageType int, msgData []byte, conn *websocket.Conn) {
	return func(messageType int, msgData []byte, conn *websocket.Conn) {
		mu.Lock()
		defer mu.Unlock()
		var msg Msg
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
		switch msg.Type {
		case TypeTextMsg:
			handleTypeTextMsg(msg, db, &msgJson)
		case TypeReadMsg:
			handleTypeReadMsg(msg, db, &msgJson)
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

// handleTypeTextMsg processing a message sent to the chat room.
// The message is saved to the database, then parsed and sent back to the client.
func handleTypeTextMsg(msg Msg, db *database.Database, msgJson *string) {
	if msg.Msg["Text"] == "" {
		return
	}
	newMsgData := map[string]interface{}{
		"user":    msg.Uid,
		"chat":    msg.ChatId,
		"text":    msg.Msg["Text"],
		"date":    time.Now(),
		"is_read": false,
	}
	_, err := db.SyncQ().Insert("chat_msg", newMsgData)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
	newMsgNotification(msg.ChatId, db)
	newMsg, err := getNewMsg(db, newMsgData)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
	err = IncrementChatMsgCountFromDb(msg.ChatId, msg.Uid, db)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
	db.AsyncQ().Wait()
	inc, _ := db.AsyncQ().LoadAsyncRes("inc")
	if inc.Error != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, inc.Error.Error())
		return
	}
	err = globalIncrementMessage(&newMsg, db)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
	err = setNotificationUsers(&newMsg, &msg, db)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
	*msgJson, err = newMsgJson(msg.Type, msg.Uid, msg.ChatId, newMsg)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
}

// handleTypeReadMsg processing a message read by a user.
// Changes in the database of message status to "read".
// Sending data about the message to the client.
func handleTypeReadMsg(msg Msg, db *database.Database, msgJson *string) {
	db.AsyncQ().AsyncUpdate("updMsg", "chat_msg", []dbutils.DbEquals{
		{
			Name:  "is_read",
			Value: true,
		},
	}, dbutils.WHEquals(map[string]interface{}{"id": msg.Msg["Id"]}, "AND"))
	db.AsyncQ().AsyncQuery("decMsgCount", "UPDATE `chat_msg_count` SET `count`= `count` - 1 WHERE user = ? AND chat = ? ;", msg.Uid, msg.ChatId)
	db.AsyncQ().Wait()
	updMsg, _ := db.AsyncQ().LoadAsyncRes("updMsg")
	if updMsg.Error != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, updMsg.Error.Error())
		return
	}
	decMsgCount, _ := db.AsyncQ().LoadAsyncRes("decMsgCount")
	if decMsgCount.Error != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, decMsgCount.Error.Error())
		return
	}
	var err error
	*msgJson, err = newMsgJson(msg.Type, msg.Uid, msg.ChatId, msg.Msg)
	if err != nil {
		*msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		return
	}
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

func setNotificationUsers(newMsg *map[string]string, msg *Msg, db *database.Database) error {
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"id": msg.ChatId,
	}, "AND"), 1)
	if err != nil {
		return err
	}
	uidInt, err := strconv.Atoi(msg.Uid)
	if err != nil {
		return err
	}
	user1, err := dbutils.ParseInt(chat[0]["user1"])
	if err != nil {
		return err
	}
	user2, err := dbutils.ParseInt(chat[0]["user2"])
	if err != nil {
		return err
	}
	if user1 != uidInt {
		(*newMsg)["SendToUsersId"] = strconv.Itoa(user1)
	}
	if user2 != uidInt {
		(*newMsg)["SendToUsersId"] = strconv.Itoa(user2)
	}
	return nil
}

func newMsgNotification(chatId string, db *database.Database) {
	db.AsyncQ().AsyncCount("notification", []string{"id"}, "chat_msg", dbutils.WHEquals(map[string]interface{}{
		"chat":    chatId,
		"is_read": false,
	}, "AND"), 1)
}

func globalIncrementMessage(msg *map[string]string, db *database.Database) error {
	res, ok := db.AsyncQ().LoadAsyncRes("notification")
	(*msg)["GlobalIncrement"] = "1"
	if !ok {
		return nil
	}
	if res.Error != nil {
		return res.Error
	}
	count, err := dbutils.ParseInt(res.Res[0]["COUNT(id)"])
	if err != nil {
		return err
	}
	if count == 1 {
		(*msg)["GlobalIncrement"] = "0"
	}
	return nil
}

func IncrementChatMsgCountFromDb(chatId string, sendUid string, db *database.Database) error {
	recipientUid, err := getRecipientUid(chatId, sendUid, db)
	if err != nil {
		return err
	}
	db.AsyncQ().AsyncQuery("inc", "UPDATE `chat_msg_count` SET `count`= `count` + 1 WHERE user = ? AND chat = ? ;", recipientUid, chatId)
	return nil
}

func getRecipientUid(chatId string, sendUid string, db *database.Database) (int, error) {
	chat, err := db.SyncQ().QB().Select("*", "chat").
		Where("id", "=", chatId, "AND", "user1", "=", sendUid, "OR", "user2", "=", sendUid).Ex()
	if err != nil {
		return 0, err
	}
	sendUidInt, err := strconv.Atoi(sendUid)
	if err != nil {
		return 0, err
	}
	user1, err := dbutils.ParseInt(chat[0]["user1"])
	if err != nil {
		return 0, err
	}
	user2, err := dbutils.ParseInt(chat[0]["user2"])
	if err != nil {
		return 0, err
	}
	var recipientId int
	if user1 == sendUidInt {
		recipientId = user2
	}
	if user2 == sendUidInt {
		recipientId = user1
	}
	return recipientId, nil
}

// wsError sending the error to the client.
func wsError(uid string, chatId string, error string) string {
	msgJson, err := newMsgJson(TypeError, uid, chatId, map[string]string{"Error": error})
	if err != nil {
		panic(err)
	}
	return msgJson
}

// newMsgJson sending a response to the client in json format.
func newMsgJson(_type int, uid string, chatId string, msg map[string]string) (string, error) {
	m := Msg{
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
