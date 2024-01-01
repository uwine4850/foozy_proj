package chat

import (
	"encoding/json"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/notification"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
	"strconv"
	"time"
)

type chatForm struct {
	Id      []string `form:"chatId"`
	MsgText []string `form:"msg"`
	UserId  []string `form:"userId"`
}

func CreateChatPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	// Parse form.
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	var chatF chatForm
	err = form.FillStructFromForm(newForm, form.NewFillableFormStruct(&chatF), []string{})
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	if chatF.Id[0] != "-1" {
		return func() { http.Redirect(w, r, "/chat/"+chatF.Id[0], http.StatusFound) }
	}

	// Create chat
	uid, err := r.Cookie("UID")
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	db := conf.NewDb()
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

	// Creating a chat room in the database.
	_, err = db.SyncQ().Insert("chat", map[string]interface{}{"user1": uid.Value, "user2": chatF.UserId[0]})
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	// Searching for the id of a new chat room.
	chat, err := db.SyncQ().Select([]string{"id"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"user1": uid.Value,
		"user2": chatF.UserId[0],
	}, "AND"), 1)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	// If a new chat exists, creating the first message and redirecting to the chat page.
	if chat != nil {
		chatId, err := dbutils.ParseInt(chat[0]["id"])
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		_, err = db.SyncQ().Insert("chat_msg", map[string]interface{}{
			"user":    uid.Value,
			"chat":    chatId,
			"text":    chatF.MsgText[0],
			"date":    time.Now(),
			"is_read": false,
		})
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		err = insertMsgCount(chatId, uid.Value, chatF.UserId[0], db)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		err = newChat(r, strconv.Itoa(chatId), uid.Value, chatF.UserId[0], chatF.MsgText[0], db)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		err = IncrementChatMsgCountFromDb(r, strconv.Itoa(chatId), uid.Value, chatF.MsgText[0], db)
		if err != nil {
			router.ServerError(w, err.Error())
		}
		return func() { http.Redirect(w, r, "/chat/"+strconv.Itoa(chatId), http.StatusFound) }
	}
	return func() {}
}

func insertMsgCount(chatId int, uid1 string, uid2 string, db *database.Database) error {
	insertData := map[string]interface{}{
		"chat": chatId,
	}
	insertData["user"] = uid1
	_, err := db.SyncQ().Insert("chat_msg_count", insertData)
	if err != nil {
		return err
	}
	insertData["user"] = uid2
	_, err = db.SyncQ().Insert("chat_msg_count", insertData)
	if err != nil {
		return err
	}
	return nil
}

func newChat(r *http.Request, chatId string, createUID string, uid string, text string, db *database.Database) error {
	messageData := make(map[string]string)
	userDataById, err := profile.GetUserDataById(createUID, db)
	if err != nil {
		return err
	}
	userDataByIdJson, err := json.Marshal(userDataById)
	if err != nil {
		return err
	}
	messageData["User"] = string(userDataByIdJson)
	messageData["Text"] = text
	messageData["ChatId"] = chatId
	err = notification.SendNewChat(r, uid, &messageData)
	if err != nil {
		return err
	}
	err = notification.SendGlobalIncrementMsg(r, uid)
	if err != nil {
		return err
	}
	return nil
}
