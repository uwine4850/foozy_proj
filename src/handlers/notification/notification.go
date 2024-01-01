package notification

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
)

const (
	WsConnect = iota
	WsError
	WsIncrementChatMsgCount
	WsGlobalIncrementMsg
	WsGlobalDecrementMsg
	WsPopUpMessage
)

type Notification struct {
	Type    int
	UserIds []string
	Msg     map[string]string
}

type ActionFunc func(messageJsonData *[]byte, notificationData *Notification, conn *websocket.Conn)

var actionsMap = map[int]ActionFunc{
	WsIncrementChatMsgCount: handleWsIncrementChatMsgCount,
	WsGlobalIncrementMsg:    handleWsGlobalIncrementMsg,
	WsGlobalDecrementMsg:    handleWsGlobalDecrementMsg,
	WsPopUpMessage:          handleWsPopUpMessage,
}

var connections = make(map[*websocket.Conn]string)

func WsHandler(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	ws := manager.CurrentWebsocket()
	var notificationJsonData []byte
	ws.OnConnect(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		once := GetRequestOnce(r)
		uidCookie, _ := r.Cookie("UID")
		if once == "false" {
			connections[conn] = uidCookie.Value
		}
		err := ws.SendMessage(websocket.TextMessage, notificationJsonData, conn)
		if err != nil {
			panic(err)
		}
	})
	ws.OnClientClose(func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
		once := GetRequestOnce(r)
		if once == "false" {
			delete(connections, conn)
		} else {
			once = "false"
		}
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	})
	ws.OnMessage(func(messageType int, msgData []byte, conn *websocket.Conn) {
		var messageJsonNotification []byte
		var notification Notification
		isError := false
		err := json.Unmarshal(msgData, &notification)
		if err != nil {
			messageJsonNotification, _ = notificationError(err)
			isError = true
		}

		if !isError {
			actionFunc, ok := actionsMap[notification.Type]
			if ok {
				actionFunc(&messageJsonNotification, &notification, conn)
			}
		}
		for i := 0; i < len(notification.UserIds); i++ {
			for key, value := range connections {
				if notification.UserIds[i] == value {
					err := ws.SendMessage(messageType, messageJsonNotification, key)
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

func handleWsIncrementChatMsgCount(messageJsonData *[]byte, notificationData *Notification, conn *websocket.Conn) {
	nj, err := notificationJson(notificationData.Type, notificationData.UserIds, notificationData.Msg)
	if err != nil {
		njError, err := notificationError(err)
		if err != nil {
			panic(err)
		}
		*messageJsonData = njError
		return
	}
	*messageJsonData = nj
}

func handleWsGlobalIncrementMsg(messageJsonData *[]byte, notificationData *Notification, conn *websocket.Conn) {
	nj, err := notificationJson(notificationData.Type, notificationData.UserIds, notificationData.Msg)
	if err != nil {
		njError, _ := notificationError(err)
		*messageJsonData = njError
		return
	}
	*messageJsonData = nj
}

func handleWsGlobalDecrementMsg(messageJsonData *[]byte, notificationData *Notification, conn *websocket.Conn) {
	nj, err := notificationJson(notificationData.Type, notificationData.UserIds, notificationData.Msg)
	if err != nil {
		njError, _ := notificationError(err)
		*messageJsonData = njError
		return
	}
	*messageJsonData = nj
}

func handleWsPopUpMessage(messageJsonData *[]byte, notificationData *Notification, conn *websocket.Conn) {
	nj, err := notificationJson(notificationData.Type, notificationData.UserIds, notificationData.Msg)
	if err != nil {
		njError, _ := notificationError(err)
		*messageJsonData = njError
		return
	}
	*messageJsonData = nj
}

func notificationJson(_type int, usersIds []string, msg map[string]string) ([]byte, error) {
	n := Notification{
		Type:    _type,
		UserIds: usersIds,
		Msg:     msg,
	}
	marshal, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func notificationError(err error) ([]byte, error) {
	return notificationJson(WsError, []string{}, map[string]string{"error": err.Error()})
}

// GetRequestOnce gets the result of the "once" data from the header.
func GetRequestOnce(r *http.Request) string {
	once := r.Header.Get("once")
	if once == "" {
		return "false"
	}
	r.Header.Del("once")
	return once
}

func SendIncrementMsgChatCount(r *http.Request, userId string, chatId string) error {
	nj, err := notificationJson(WsIncrementChatMsgCount, []string{userId}, map[string]string{"chatId": chatId})
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, string(nj), "ws://localhost:8000/notification-ws", true)
	if err != nil {
		return err
	}
	return nil
}

func SendGlobalIncrementMsg(r *http.Request, userId string) error {
	nj, err := notificationJson(WsGlobalIncrementMsg, []string{userId}, map[string]string{})
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, string(nj), "ws://localhost:8000/notification-ws", true)
	if err != nil {
		return err
	}
	return nil
}

func SendGlobalDecrementMsg(r *http.Request, userId string) error {
	nj, err := notificationJson(WsGlobalDecrementMsg, []string{userId}, map[string]string{})
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, string(nj), "ws://localhost:8000/notification-ws", true)
	if err != nil {
		return err
	}
	return nil
}

func SendPopUpMessage(r *http.Request, userId string, messageData *map[string]string) error {
	nj, err := notificationJson(WsPopUpMessage, []string{userId}, *messageData)
	if err != nil {
		return err
	}
	err = utils.WsSendMessage(r, string(nj), "ws://localhost:8000/notification-ws", true)
	if err != nil {
		return err
	}
	return nil
}
