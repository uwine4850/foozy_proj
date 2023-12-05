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

// getChatData Returns data from the database about a specific chat.
func getChatData(id int, db *database.Database) (map[string]interface{}, error) {
	chat, err := db.SyncQ().Select([]string{"*"}, "chat", dbutils.WHEquals(map[string]interface{}{
		"id": id,
	}, "AND"), 1)
	if err != nil {
		return nil, err
	}
	return chat[0], nil
}

// getChatUser Returns user data from the database.
func getChatUser(chatDb map[string]interface{}, uid int, db *database.Database) (map[string]interface{}, error) {
	user1, err := dbutils.ParseInt(chatDb["user1"])
	if err != nil {
		return nil, err
	}
	user2, err := dbutils.ParseInt(chatDb["user2"])
	if err != nil {
		return nil, err
	}
	var dbUid int
	if user1 == uid {
		dbUid = user2
	}
	if user2 == uid {
		dbUid = user1
	}
	dbUser, err := db.SyncQ().Select([]string{"*"}, "auth", dbutils.WHEquals(map[string]interface{}{
		"id": dbUid,
	}, "AND"), 1)
	if err != nil {
		return nil, err
	}
	delete(dbUser[0], "password")
	return dbUser[0], nil
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
	uidInt, _ := strconv.Atoi(uid.Value)

	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		panic(err)
		//return func() { router.ServerError(w, err.Error()) }
	}
	defer db.Close()

	chatDb, err := getChatData(chatIdInt, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	user, err := getChatUser(chatDb, uidInt, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	var userData profile.UserData
	err = dbutils.FillStructFromDb(user, &userData)
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
