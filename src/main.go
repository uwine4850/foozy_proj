package main

import (
	"github.com/uwine4850/foozy/pkg/builtin/builtin_mddl"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/middlewares"
	"github.com/uwine4850/foozy/pkg/router"
	server2 "github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/tmlengine"
	"github.com/uwine4850/foozy_proj/src/handlers"
	"github.com/uwine4850/foozy_proj/src/middlewares/profilemddl"
	"net/http"
)

func main() {
	mddl := middlewares.NewMiddleware()
	mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
	mddl.HandlerMddl(1, profilemddl.AuthMddl)
	engine, err := tmlengine.NewTemplateEngine()
	if err != nil {
		panic(err)
	}
	manager := router.NewManager(engine)
	newRouter := router.NewRouter(manager)
	newRouter.EnableLog(true)
	newRouter.SetMiddleware(mddl)
	newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/home.html")
		UID, _ := manager.GetUserContext("UID")
		manager.SetContext(map[string]interface{}{"UID": UID.(string)})
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/prof/<id>", handlers.ProfileView)
	newRouter.Get("/new-post", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/new_post.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/register", handlers.Register)
	newRouter.Post("/register-post", handlers.RegisterPost)
	newRouter.Get("/sign-in", handlers.SignIn)
	newRouter.Post("/sign-in-post", handlers.SignInPost)
	newRouter.Get("/profile/<id>/edit", handlers.ProfileEdit)
	newRouter.Post("/profile-edit-post/<id>", handlers.ProfileEditPost)
	newRouter.Post("/log-out-post", handlers.ProfileLogOutPost)
	newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	newRouter.GetMux().Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))
	server := server2.NewServer(":8000", newRouter)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
