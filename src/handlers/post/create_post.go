package post

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/handlers/profile"
	"net/http"
	"time"
)

type NewPostFill struct {
	Name        []string        `form:"name"`
	Description []string        `form:"description"`
	Images      []form.FormFile `form:"images"`
	Category    []string        `form:"category"`
}

type Post struct {
	Id           string `db:"id"`
	ParentUserId string `db:"parent_user"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	CategoryId   string `db:"category"`
	Date         string `db:"date"`
}

type Category struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func CreatePost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	db := conf.DatabaseI
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

	categories, err := db.SyncQ().QB().Select("*", "post_categories").Ex()
	if err != nil {
		panic(err)
	}
	var categoriesStruct []Category
	for i := 0; i < len(categories); i++ {
		var category Category
		err := dbutils.FillStructFromDb(categories[i], &category)
		if err != nil {
			return func() { router.ServerError(w, err.Error()) }
		}
		categoriesStruct = append(categoriesStruct, category)
	}

	manager.SetTemplatePath("src/templates/new_post.html")
	manager.SetContext(map[string]interface{}{"categories": categoriesStruct})
	router.HandleRedirectError(manager)
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
	return func() {}
}

func SavePost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	frm := form.NewForm(r)
	err := frm.Parse()
	if err != nil {
		panic(err)
	}
	var newPostFill NewPostFill
	filledPostStruct := form.NewFillableFormStruct(&newPostFill)
	err = form.FillStructFromForm(frm, filledPostStruct, []string{})
	if err != nil {
		return func() { router.RedirectError(w, r, "/new-post", err.Error(), manager) }
	}
	err = form.FieldsNotEmpty(&newPostFill, []string{"Name", "Images", "Category"})
	if err != nil {
		return func() { router.RedirectError(w, r, "/new-post", err.Error(), manager) }
	}
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		return nil
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	user, _ := manager.GetUserContext("currentUser")
	userData := user.(profile.UserData)
	_, err = db.SyncQ().Insert("posts", map[string]interface{}{
		"parent_user": userData.Id,
		"name":        filledPostStruct.GetOrDef("Name", 0),
		"description": filledPostStruct.GetOrDef("Description", 0),
		"category":    filledPostStruct.GetOrDef("Category", 0),
		"date":        time.Now(),
	})
	err = saveImages(w, &newPostFill, db)
	if err != nil {
		return func() { router.RedirectError(w, r, "/new-post", err.Error(), manager) }
	}
	if err != nil {
		return func() { router.RedirectError(w, r, "/new-post", err.Error(), manager) }
	}
	return func() { http.Redirect(w, r, "/new-post", http.StatusFound) }
}

func saveImages(w http.ResponseWriter, newPost *NewPostFill, db *database.Database) error {
	postQ, err := db.SyncQ().Query("SELECT * FROM `posts` ORDER BY date DESC LIMIT 1;")
	if err != nil {
		return err
	}
	if postQ == nil {
		return errors.New("created post not found")
	}
	var post Post
	err = dbutils.FillStructFromDb(postQ[0], &post)
	if err != nil {
		return err
	}
	for i := 0; i < len(newPost.Images); i++ {
		newPostFormFile := newPost.Images[i]
		var buildPath string
		err = form.SaveFile(w, newPostFormFile.Header, "media/posts", &buildPath)
		if err != nil {
			return err
		}
		_, err = db.SyncQ().Insert("post_images", map[string]interface{}{
			"postid": post.Id,
			"path":   buildPath,
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}
