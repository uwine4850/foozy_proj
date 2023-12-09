package chat

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
)

type chat struct {
	Id    string `db:"id"`
	User1 string `db:"user1"`
	User2 string `db:"user2"`
}

type chatInfo struct {
	Chat     chat
	User     profile.UserData
	LastMsg  ChatMessage
	MsgCount int
}

func ChatList(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	uid, err := r.Cookie("UID")
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	//Database Connection.
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)

	// Retrieving chats associated with a user.
	chats, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"user1": uid.Value,
		"user2": uid.Value,
	}, "OR"), 0)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	// Set chat information.
	var chatsStruct []chat
	var chatsInfo []chatInfo
	for i := 0; i < len(chats); i++ {
		var c chat
		err := dbutils.FillStructFromDb(chats[i], &c)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		chatsInfo = append(chatsInfo, chatInfo{
			Chat: c,
		})
		chatsStruct = append(chatsStruct, c)
	}
	// Set user information.
	err = getChatListUsers(chatsStruct, uid.Value, &chatsInfo, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	// Set information about the latest messages.
	err = setChatLastMsg(chatsStruct, &chatsInfo, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	err = setMsgCount(uid.Value, chatsStruct, &chatsInfo, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetContext(map[string]interface{}{"chatsInfo": chatsInfo})
	manager.SetTemplatePath("src/templates/chat_list.html")
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

// getChatListUsers Sets the user information for each instance of the chatInfo structure.
func getChatListUsers(chats []chat, uid string, _chatInfo *[]chatInfo, db *database.Database) error {
	var asyncKeys []string
	for i := 0; i < len(chats); i++ {
		var currentChatUserId string
		if chats[i].User1 == uid {
			currentChatUserId = chats[i].User2
		} else {
			currentChatUserId = chats[i].User1
		}
		db.AsyncQ().AsyncSelect("user"+chats[i].Id, []string{"*"}, "auth", dbutils.WHEquals(map[string]interface{}{
			"id": currentChatUserId,
		}, "AND"), 1)
		asyncKeys = append(asyncKeys, "user"+chats[i].Id)
	}
	db.AsyncQ().Wait()
	for i := 0; i < len(asyncKeys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(asyncKeys[i])
		if !ok {
			return errors.New("database output error: the required LoadAsyncRes result was not found")
		}
		if res.Error != nil {
			return res.Error
		}
		var userData profile.UserData
		err := dbutils.FillStructFromDb(res.Res[0], &userData)
		if err != nil {
			return err
		}
		(*_chatInfo)[i].User = userData
	}
	return nil
}

// setChatLastMsg Sets the last message in the chatInfo structure.
// For each chat, an asynchronous request is sent to the database to retrieve the last message.
func setChatLastMsg(chats []chat, _chatInfo *[]chatInfo, db *database.Database) error {
	var asyncKeys []string
	for i := 0; i < len(chats); i++ {
		db.AsyncQ().AsyncSelect("msg"+chats[i].Id, []string{"*"}, "chat_msg", dbutils.WHOutput{
			QueryStr:  "chat = ? ORDER BY id DESC",
			QueryArgs: []interface{}{chats[i].Id},
		}, 1)
		asyncKeys = append(asyncKeys, "msg"+chats[i].Id)
	}
	db.AsyncQ().Wait()
	for i := 0; i < len(asyncKeys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(asyncKeys[i])
		if !ok {
			return errors.New("database output error: the required LoadAsyncRes result was not found")
		}
		if res.Error != nil {
			return res.Error
		}
		var chatMsg ChatMessage
		if res.Res != nil {
			err := dbutils.FillStructFromDb(res.Res[0], &chatMsg)
			if err != nil {
				return err
			}
		}
		(*_chatInfo)[i].LastMsg = chatMsg
	}
	return nil
}

func setMsgCount(uid string, chats []chat, _chatInfo *[]chatInfo, db *database.Database) error {
	var asyncKeys []string
	for i := 0; i < len(chats); i++ {
		db.AsyncQ().AsyncSelect("count"+chats[i].Id, []string{"count"}, "chat_msg_count", dbutils.WHEquals(map[string]interface{}{
			"chat": chats[i].Id,
			"user": uid,
		}, "AND"), 1)
		asyncKeys = append(asyncKeys, "count"+chats[i].Id)
	}
	db.AsyncQ().Wait()
	for i := 0; i < len(asyncKeys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(asyncKeys[i])
		if !ok {
			return errors.New("database output error: the result of the number of messages query was not found")
		}
		if res.Error != nil {
			return res.Error
		}
		count, err := dbutils.ParseInt(res.Res[0]["count"])
		if err != nil {
			return err
		}
		(*_chatInfo)[i].MsgCount = count
	}
	return nil
}
