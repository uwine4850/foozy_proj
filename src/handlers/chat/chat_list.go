package chat

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"net/http"
)

func ChatList(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	manager.SetTemplatePath("src/templates/chat_list.html")
	err := manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}
