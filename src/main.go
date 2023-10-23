package main

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	server2 "github.com/uwine4850/foozy/pkg/server"
	"github.com/uwine4850/foozy/pkg/tmlengine"
	"net/http"
)

func main() {
	engine, err := tmlengine.NewTemplateEngine()
	if err != nil {
		panic(err)
	}
	manager := router.NewManager(engine)
	newRouter := router.NewRouter(manager)
	newRouter.EnableLog(true)
	newRouter.Get("/home", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/home.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/profile", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/profile.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/new-post", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/new_post.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/register", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/auth/register.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.Get("/sign-in", func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
		manager.SetTemplatePath("src/templates/auth/signin.html")
		err := manager.RenderTemplate(w, r)
		if err != nil {
			panic(err)
		}
	})
	newRouter.GetMux().Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
	server := server2.NewServer(":8000", newRouter)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
