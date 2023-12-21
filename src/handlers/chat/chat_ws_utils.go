package chat

import (
	"github.com/uwine4850/foozy/pkg/database"
	"strconv"
	"strings"
)

func saveMessageImages(imagesPaths string, newMessageId string, db *database.Database) error {
	images := strings.Split(imagesPaths, "\\")
	queryKeys := make([]string, 0)
	for i := 0; i < len(images); i++ {
		db.AsyncQ().AsyncInsert("saveImg"+strconv.Itoa(i), "chat_msg_images", map[string]interface{}{
			"parent_msg": newMessageId,
			"path":       images[i],
		})
		queryKeys = append(queryKeys, "saveImg"+strconv.Itoa(i))
	}
	db.AsyncQ().Wait()
	for i := 0; i < len(queryKeys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(queryKeys[i])
		if ok {
			if res.Error != nil {
				return res.Error
			}
		}
	}
	return nil
}
