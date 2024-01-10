package chat

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy_proj/src/conf"
	"net/http"
)

func LoadImages(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	detailChatId, ok := manager.GetUserContext("detailChatId")
	if !ok {
		return func() { sendJson(map[string]interface{}{"error": "Detail Chat id not found."}, w) }
	}
	imageId := r.URL.Query().Get("imageid")
	db := conf.NewDb()
	err := db.Connect()
	if err != nil {
		return func() { sendJson(map[string]interface{}{"error": err.Error()}, w) }
	}

	images, err := LoadChatImages(detailChatId.(string), imageId, 10, db)
	if err != nil {
		return func() { sendJson(map[string]interface{}{"error": err.Error()}, w) }
	}

	// Close database.
	err = db.Close()
	if err != nil {
		return func() { sendJson(map[string]interface{}{"error": err.Error()}, w) }
	}
	return func() {
		sendJson(map[string]interface{}{"images": images}, w)
	}
}
