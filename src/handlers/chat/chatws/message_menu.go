package chatws

import (
	"encoding/json"
	"errors"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/form"
	"net/http"
	"strings"
)

func MessageMenu(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	newForm := form.NewForm(r)
	err := newForm.Parse()
	if err != nil {
		return func() { sendError(w, err) }
	}
	UID, _ := manager.GetUserContext("UID")
	action := newForm.Value("action")
	switch action {
	case "delete":
		err := SendDeleteMessage(r, UID.(string), newForm.Value("chatId"), newForm.Value("msg-id"))
		if err != nil {
			return func() { sendError(w, err) }
		}
	case "update":
		msgText := newForm.Value("messageText")
		msgId := newForm.Value("updMsgId")
		var delImages strings.Builder
		for key, _ := range newForm.GetApplicationForm() {
			if strings.HasPrefix(key, "updRmImage.") {
				path, _ := strings.CutPrefix(key, "updRmImage./")
				delImages.WriteString(path + "\\")
			}
		}
		if delImages.String() == "" && msgText == "" {
			return func() { sendError(w, errors.New("the updated message cannot be empty")) }
		}
		messageData := map[string]string{"id": msgId, "text": msgText, "delImages": delImages.String()}
		err := SendUpdateMessage(r, UID.(string), newForm.Value("chatId"), messageData)
		if err != nil {
			return func() { sendError(w, err) }
		}
	}
	return func() {}
}

func sendError(w http.ResponseWriter, err error) {
	jsonUsers, err := json.Marshal(map[string]interface{}{"error": err.Error()})
	if err != nil {
		panic(err)
	}
	w.Write(jsonUsers)
}
