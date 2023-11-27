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
	"net/http"
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

func removeElement[T comparable](slice []T, element T) []T {
	var result []T
	for _, el := range slice {
		if el != element {
			result = append(result, el)
		}
	}
	return result
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
		chatConnections[connChatId] = removeElement(chatConnections[connChatId], conn)
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	})
	ws.OnConnect(func(conn *websocket.Conn) {
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
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		mu.Lock()
		defer mu.Unlock()
		var msg Msg
		var msgJson string
		err := json.Unmarshal(msgData, &msg)
		if err != nil {
			panic(err)
		}
		db := conf.DatabaseI
		err = db.Connect()
		if err != nil {
			msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
		}
		defer func(db *database.Database) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)
		switch msg.Type {
		case TypeTextMsg:
			if msg.Msg["Text"] == "" {
				return
			}
			//db := conf.DatabaseI
			//err := db.Connect()
			//if err != nil {
			//	msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			//	break
			//}
			newMsgData := map[string]interface{}{
				"user":    msg.Uid,
				"chat":    msg.ChatId,
				"text":    msg.Msg["Text"],
				"date":    time.Now(),
				"is_read": false,
			}
			_, err = db.SyncQ().Insert("chat_msg", newMsgData)
			if err != nil {
				msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
				break
			}
			newMsg, err := getNewMsg(db, newMsgData)
			if err != nil {
				msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
				break
			}
			//err = db.Close()
			//if err != nil {
			//	msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
			//	break
			//}
			msgJson, err = newMsgJson(msg.Type, msg.Uid, msg.ChatId, newMsg)
			if err != nil {
				msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
				break
			}
		case TypeReadMsg:
			_, err := db.SyncQ().Update("chat_msg", []dbutils.DbEquals{
				{
					Name:  "is_read",
					Value: true,
				},
			}, dbutils.WHEquals(map[string]interface{}{"id": msg.Msg["Id"]}, "AND"))
			if err != nil {
				msgJson = wsError(msg.Uid, msg.ChatId, err.Error())
				break
			}
			msgJson, err = newMsgJson(msg.Type, msg.Uid, msg.ChatId, msg.Msg)
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < len(chatConnections[msg.ChatId]); i++ {
			if msgJson == "" {
				return
			}
			err = ws.SendMessage(websocket.TextMessage, []byte(msgJson), chatConnections[msg.ChatId][i])
			if err != nil {
				panic(err)
			}
		}
	})
	err := ws.ReceiveMessages(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func getNewMsg(db interfaces.IDatabase, msgData map[string]interface{}) (map[string]string, error) {
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

func wsError(uid string, chatId string, error string) string {
	msgJson, err := newMsgJson(TypeError, uid, chatId, map[string]string{"Error": error})
	if err != nil {
		panic(err)
	}
	return msgJson
}
