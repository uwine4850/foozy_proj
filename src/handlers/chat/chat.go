package chat

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
	"strconv"
)

type ChatMessage struct {
	Id     string `db:"id"`
	UserId string `db:"user"`
	Text   string `db:"text"`
	Date   string `db:"date"`
	IsRead string `db:"is_read"`
}

func Chat(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	chatId, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, "Slug parameter id for chat was not found.") }
	}
	chatIdInt, err := strconv.Atoi(chatId)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	uid, err := r.Cookie("UID")
	if err != nil {
		panic(err)
	}

	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		panic(err)
		//return func() { router.ServerError(w, err.Error()) }
	}
	defer db.Close()

	userData, err := GetRecipientUser(chatId, uid.Value, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	msg, err := loadChatMsg(chatIdInt, userData, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	var chatMessages []ChatMessage
	var cm ChatMessage
	if msg != nil {
		err = dbutils.FillStructFromDb(msg, &cm)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		chatMessages = append(chatMessages, cm)
	}
	manager.SetTemplatePath("src/templates/chat.html")
	manager.SetContext(map[string]interface{}{"user": userData, "messages": chatMessages, "chatId": chatId})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

type ErrChatNotFound struct {
}

func (e ErrChatNotFound) Error() string {
	return "chat not found"
}

func GetRecipientUser(chatId string, sendUid string, db *database.Database) (profile.UserData, error) {
	chat, err := db.SyncQ().QB().Select("*", "chat").
		Where("id", "=", chatId, "AND", "user1", "=", sendUid, "OR", "user2", "=", sendUid).Ex()
	if chat == nil {
		return profile.UserData{}, ErrChatNotFound{}
	}
	if err != nil {
		return profile.UserData{}, err
	}
	sendUidInt, err := strconv.Atoi(sendUid)
	if err != nil {
		return profile.UserData{}, err
	}
	user1, err := dbutils.ParseInt(chat[0]["user1"])
	if err != nil {
		return profile.UserData{}, err
	}
	user2, err := dbutils.ParseInt(chat[0]["user2"])
	if err != nil {
		return profile.UserData{}, err
	}
	var recipientId int
	if user1 == sendUidInt {
		recipientId = user2
	}
	if user2 == sendUidInt {
		recipientId = user1
	}
	user, err := db.SyncQ().QB().Select("*", "auth").Where("id", "=", recipientId).Ex()
	if err != nil {
		return profile.UserData{}, err
	}
	var recipientUser profile.UserData
	err = dbutils.FillStructFromDb(user[0], &recipientUser)
	if err != nil {
		return profile.UserData{}, err
	}
	return recipientUser, nil
}

// loadChatMsg Loads a single message from the database.
// If there are no read messages - returns the oldest message.
// If all messages are read - returns the most recent message.
func loadChatMsg(chatId int, userData profile.UserData, db *database.Database) (map[string]interface{}, error) {
	notReadMessage, err := db.SyncQ().Select([]string{"*"}, "chat_msg", dbutils.WHEquals(map[string]interface{}{
		"user":    userData.Id,
		"chat":    chatId,
		"is_read": 0,
	}, "AND"), 1)
	if err != nil {
		return nil, err
	}
	if notReadMessage == nil {
		message, err := db.SyncQ().Query("SELECT * FROM (SELECT * FROM `chat_msg` WHERE chat = ? "+
			"ORDER BY id DESC LIMIT 1) AS f ORDER BY id ASC;", chatId)
		if err != nil {
			return nil, err
		}
		if message == nil {
			return nil, nil
		}
		return message[0], nil
	} else {
		if notReadMessage == nil {
			return nil, nil
		}
		return notReadMessage[0], nil
	}
}
