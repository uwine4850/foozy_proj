package chat

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/form"
	"net/http"
	"strings"
)

type FormMessage struct {
	ChatId []string        `form:"chatId"`
	Uid    []string        `form:"uid"`
	Images []form.FormFile `form:"images"`
	Text   []string        `form:"text"`
}

func ReceiveMessage(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		panic(err)
	}
	var formMessage FormMessage
	fillableMessage := form.NewFillableFormStruct(&formMessage)
	err = form.FillStructFromForm(frm, fillableMessage, []string{})
	if err != nil {
		panic(err)
	}
	var message Message
	if fillableMessage.GetOrDef("Text", 0) != "" && formMessage.Images == nil {
		message = Message{
			Type:   0,
			Uid:    formMessage.Uid[0],
			ChatId: formMessage.ChatId[0],
			Msg:    map[string]string{"Text": formMessage.Text[0]},
		}
		err = SendTextMessage(r, &message)
		if err != nil {
			panic(err)
		}
	}
	if formMessage.Images != nil {
		imagesPaths, err := saveImages(w, &formMessage.Images)
		if err != nil {
			panic(err)
		}
		message = Message{
			Type:   0,
			Uid:    formMessage.Uid[0],
			ChatId: formMessage.ChatId[0],
			Msg:    map[string]string{"Text": fillableMessage.GetOrDef("Text", 0), "images": imagesPaths},
		}
		err = SendImageMessage(r, &message)
		if err != nil {
			panic(err)
		}
	}
	return func() {}
}

func saveImages(w http.ResponseWriter, images *[]form.FormFile) (string, error) {
	var paths = make([]string, 0)
	for i := 0; i < len(*images); i++ {
		var path string
		err := form.SaveFile(w, (*images)[i].Header, "media/chat_images", &path)
		if err != nil {
			return "", err
		}
		paths = append(paths, path)
	}
	res := strings.Join(paths, "\\")
	return res, nil
}
