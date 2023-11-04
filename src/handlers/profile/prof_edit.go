package profile

import (
	"errors"
	"fmt"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/form"
	utils2 "github.com/uwine4850/foozy/pkg/utils"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"net/http"
	"os"
	"strconv"
)

func ProfileEdit(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if !editPermission(w, r, manager) {
		return func() {}
	}
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
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id")) }
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	user, err := db.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", uidInt}}, 1)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	var fillUserData UserData
	err = dbutils.FillStructFromDb(user[0], &fillUserData)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	manager.SetTemplatePath("src/templates/profile_edit.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": "", "user": fillUserData})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string), "user": fillUserData})
		manager.DelUserContext("error")
	}
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {}
}

type editFormData struct {
	Name        []string `form:"name"`
	Description []string `form:"description"`
	DelAvatar   string
	AvatarPath  string
}

func ProfileEditPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if !editPermission(w, r, manager) {
		return func() {}
	}
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		return func() { router.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id")) }
	}
	// Get form data
	frm, err := utils.ParseForm(r)
	if err != nil {
		return func() { router.RedirectError(w, r, fmt.Sprintf("/profile/%s/edit", uid), err.Error(), manager) }
	}
	var fillProfileEditForm editFormData
	err = form.FillStructFromForm(frm, &fillProfileEditForm)
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	fillProfileEditForm.DelAvatar = frm.Value("del_avatar")

	// Save avatar
	if fillProfileEditForm.DelAvatar == "" {
		file, fileHeader, err := frm.File("avatar")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			return func() { router.ServerError(w, err.Error()) }
		}
		var buildPath string
		if !errors.Is(err, http.ErrMissingFile) {
			err := form.SaveFile(w, file, fileHeader, "media/avatars/", &buildPath)
			if err != nil {
				return func() { router.ServerError(w, err.Error()) }
			}
		}
		fillProfileEditForm.AvatarPath = buildPath
	}

	// Update profile
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(db)
	updSlice := []dbutils.DbEquals{
		{"name", fillProfileEditForm.Name[0]},
		{"description", fillProfileEditForm.Description[0]},
	}
	// Delete avatar
	user, err := db.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", uid}}, 1)
	dbAvatarPath := dbutils.ParseString(user[0]["avatar"])
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	if fillProfileEditForm.DelAvatar != "" {
		updSlice = append(updSlice, dbutils.DbEquals{Name: "avatar", Value: ""})
		if utils2.PathExist(dbAvatarPath) {
			err := os.Remove(dbAvatarPath)
			if err != nil {
				return func() { router.ServerError(w, err.Error()) }
			}
		}
	} else {
		// Delete old avatar
		if fillProfileEditForm.AvatarPath != "" {
			updSlice = append(updSlice, dbutils.DbEquals{Name: "avatar", Value: fillProfileEditForm.AvatarPath})
			if utils2.PathExist(dbAvatarPath) {
				err := os.Remove(dbAvatarPath)
				if err != nil {
					return func() { router.ServerError(w, err.Error()) }
				}
			}
		}
	}
	_, err = db.SyncQ().Update("auth", updSlice, []dbutils.DbEquals{{"id", uid}})
	if err != nil {
		return func() { router.ServerError(w, err.Error()) }
	}
	return func() {
		http.Redirect(w, r, fmt.Sprintf("/prof/%s", uid), http.StatusFound)
	}
}

func editPermission(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) bool {
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		router.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id"))
		return false
	}
	uidC, err := r.Cookie("UID")
	if err != nil {
		router.ServerError(w, err.Error())
		return false
	}
	if uid != uidC.Value {
		router.ServerForbidden(w)
		return false
	}
	return true
}
