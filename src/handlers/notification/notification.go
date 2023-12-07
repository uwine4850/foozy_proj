package notification

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
	"sync"
)

const (
	TypeConnect = iota
	TypeGlobalIncrementMsg
)

type Notification struct {
	Type int
	Uid  []string
	Msg  map[string]string
}

var uidConnections = make(map[string]*websocket.Conn)
var connections = make(map[*websocket.Conn]string)

func NotificationWs(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	var mu sync.Mutex
	ws := manager.GetWebSocket()
	ws.OnConnect(func(conn *websocket.Conn) {
		uid, err := r.Cookie("UID")
		if err != nil {
			panic(err)
		}
		uidConnections[uid.Value] = conn
		connections[conn] = uid.Value
		notificationJson, err := newNotificationJson(TypeConnect, []string{}, map[string]string{})
		if err != nil {
			panic(err)
		}
		err = ws.SendMessage(websocket.TextMessage, []byte(notificationJson), conn)
		if err != nil {
			panic(err)
		}
	})
	ws.OnClientClose(func(conn *websocket.Conn) {
		uid, ok := connections[conn]
		if ok {
			delete(uidConnections, uid)
			delete(connections, conn)
		}
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		mu.Lock()
		defer mu.Unlock()
		var notification Notification
		err := json.Unmarshal(msgData, &notification)
		if err != nil {
			panic(err)
		}
		var notificationJson string
		switch notification.Type {
		case TypeGlobalIncrementMsg:
			nj, err := newNotificationJson(notification.Type, notification.Uid, notification.Msg)
			if err != nil {
				panic(err)
			}
			notificationJson = nj
		}
		if notificationJson != "" {
			for i := 0; i < len(notification.Uid); i++ {
				uidConn := uidConnections[notification.Uid[i]]
				if conn != nil {
					err := ws.SendMessage(messageType, []byte(notificationJson), uidConn)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	})
	err := ws.ReceiveMessages(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func newNotificationJson(_type int, uid []string, msg map[string]string) (string, error) {
	m := Notification{
		Type: _type,
		Uid:  uid,
		Msg:  msg,
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
