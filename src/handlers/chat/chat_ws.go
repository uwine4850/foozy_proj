package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"time"
)

const (
	TypeConnect = iota
	TypeTextMsg
)

type TextMsg struct {
	Text string
}

type Msg struct {
	Type   int
	Uid    string
	ChatId string
	Msg    TextMsg
}

func newMsgJson(_type int, uid string, chatId string, msg TextMsg) (string, error) {
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
		msgJson, err := newMsgJson(TypeConnect, uid.Value, chatId, TextMsg{})
		if err != nil {
			panic(err)
		}
		err = ws.SendMessage(websocket.TextMessage, []byte(msgJson), conn)
		if err != nil {
			panic(err)
		}
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		var msg Msg
		var msgJson string
		err := json.Unmarshal(msgData, &msg)
		if err != nil {
			panic(err)
		}
		switch msg.Type {
		case TypeTextMsg:
			if msg.Msg.Text == "" {
				return
			}
			db := conf.DatabaseI
			err := db.Connect()
			if err != nil {
				panic(err)
			}
			_, err = db.SyncQ().Insert("chat_msg", map[string]interface{}{
				"user": msg.Uid,
				"chat": msg.ChatId,
				"text": msg.Msg.Text,
				"date": time.Now(),
			})
			if err != nil {
				panic(err)
			}
			err = db.Close()
			if err != nil {
				panic(err)
			}
			msgJson, err = newMsgJson(msg.Type, msg.Uid, msg.ChatId, msg.Msg)
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < len(chatConnections[msg.ChatId]); i++ {
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
