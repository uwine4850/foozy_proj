package chat

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
	"strconv"
)

func Detail(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	_uid, _ := manager.GetUserContext("UID")
	UID := _uid.(string)
	chatId, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, "Slug parameter id for Chat was not found.") }
	}
	manager.SetUserContext("detailChatId", chatId)

	db := conf.NewDb()
	err := db.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)

	chat, err := GetChat(chatId, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	// Check permission
	if !checkPermission(chat, UID) {
		return func() { router.ServerForbidden(w) }
	}

	images, err := LoadChatImages(chatId, "", 10, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	user, err := GetRecipientUser(chatId, UID, db)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}

	manager.SetTemplatePath("src/templates/chat_detail.html")
	manager.SetContext(map[string]interface{}{"images": images, "user": user})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

func LoadChatImages(chatId string, imageId string, count int, db *database.Database) ([]MessageImage, error) {
	var _imageId string
	if imageId != "" {
		_imageId = "AND chat_msg_images.id < " + imageId
	}
	imagesQuery, err := db.SyncQ().Query(fmt.Sprintf("SELECT chat_msg_images.* "+
		"FROM chat_msg_images "+
		"INNER JOIN chat_msg ON chat_msg_images.parent_msg = chat_msg.id "+
		"WHERE chat_msg.Chat = %s %s ORDER BY chat_msg_images.id DESC LIMIT %s;",
		chatId, _imageId, strconv.Itoa(count)))
	if err != nil {
		return nil, err
	}
	var images []MessageImage
	for i := 0; i < len(imagesQuery); i++ {
		var image MessageImage
		err := dbutils.FillStructFromDb(imagesQuery[i], &image)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func checkPermission(chat Chat, uid string) bool {
	if chat.User1 != uid && chat.User2 != uid {
		return false
	}
	return true
}

func SearchMessages(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		panic(err)
	}
	detailChatId, ok := manager.GetUserContext("detailChatId")
	if !ok {
		return func() { router.SendJson(map[string]interface{}{"error": "Detail Chat id not found."}, w) }
	}
	UID, _ := manager.GetUserContext("UID")
	messageText := frm.Value("message-text")

	db := conf.NewDb()
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.SendJson(map[string]interface{}{"error": err.Error()}, w)
		}
	}(db)

	messagesDb, err := db.SyncQ().QB().Select("*", "chat_msg").
		Where("chat", "=", detailChatId, "AND", "text LIKE \"%"+messageText+"%\" ORDER BY id DESC LIMIT 5").Ex()
	if err != nil {
		return func() { router.SendJson(map[string]interface{}{"error": err.Error()}, w) }
	}

	var messages []ChatMessage
	for i := 0; i < len(messagesDb); i++ {
		var m ChatMessage
		err := dbutils.FillStructFromDb(messagesDb[i], &m)
		if err != nil {
			return func() { router.SendJson(map[string]interface{}{"error": err.Error()}, w) }
		}
		images, err := LoadMessageImages(m.Id, db)
		if err != nil {
			return func() { router.SendJson(map[string]interface{}{"error": err.Error()}, w) }
		}
		m.Images = images
		messages = append(messages, m)
	}

	return func() {
		router.SendJson(map[string]interface{}{"messages": messages, "UID": UID}, w)
	}
}
