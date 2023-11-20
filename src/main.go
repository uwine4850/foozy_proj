package main

import (
	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/middlewares"
	"github.com/uwine4850/foozy/pkg/router"
	server2 "github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/tmlengine"
	"github.com/uwine4850/foozy_proj/src/handlers/chat"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"github.com/uwine4850/foozy_proj/src/middlewares/chatmddl"
	"github.com/uwine4850/foozy_proj/src/middlewares/profilemddl"
	"net/http"
)

func main() {
	mddl := middlewares.NewMiddleware()
	mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
	mddl.AsyncHandlerMddl(chatmddl.ChatPermissionMddl)
	mddl.HandlerMddl(1, profilemddl.AuthMddl)
	engine, err := tmlengine.NewTemplateEngine()
	if err != nil {
		panic(err)
	}
	manager := router.NewManager(engine)
	newRouter := router.NewRouter(manager)
	newRouter.EnableLog(true)
	newRouter.SetMiddleware(mddl)
	newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		_, ok := manager.GetUserContext("mddl_error")
		if ok {
			return func() {}
		}
		manager.SetTemplatePath("src/templates/home.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
		return func() {}
	})
	newRouter.Get("/prof/<id>", profile.ProfileView)
	newRouter.Get("/new-post", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
		manager.SetTemplatePath("src/templates/new_post.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
		return func() {}
	})
	newRouter.Get("/register", profile.Register)
	newRouter.Post("/register-post", profile.RegisterPost)
	newRouter.Get("/sign-in", profile.SignIn)
	newRouter.Post("/sign-in-post", profile.SignInPost)
	newRouter.Get("/profile/<id>/edit", profile.ProfileEdit)
	newRouter.Post("/profile-edit-post/<id>", profile.ProfileEditPost)
	newRouter.Post("/log-out-post", profile.ProfileLogOutPost)
	newRouter.Post("/subscribe-post", profile.SubscribePost)
	newRouter.Get("/my-subscriptions", profile.MySubscriptions)
	newRouter.Get("/chat/<id>", chat.Chat)
	newRouter.Get("/chat-list", chat.ChatList)
	newRouter.Post("/create-chat", chat.CreateChatPost)
	newRouter.Ws("/chat-ws", router.NewWebsocket(router.Upgrader), chat.ChatWs)
	newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	newRouter.GetMux().Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))
	server := server2.NewServer(":8000", newRouter)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
