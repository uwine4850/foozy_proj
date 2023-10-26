package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/uwine4850/foozy/pkg/builtin/auth"
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	utils2 "github.com/uwine4850/foozy/pkg/utils"
	"github.com/uwine4850/foozy_proj/src/conf"
	"github.com/uwine4850/foozy_proj/src/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func Register(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	manager.SetTemplatePath("src/templates/auth/register.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
	err := manager.RenderTemplate(w, r)
	if err != nil {
		utils.ServerError(w, err.Error())
	}
}

func RegisterPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm, err := utils.ParseForm(r)
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"name", "username", "password", "confirm_pass"}, frm.GetApplicationForm())
	if !ok {
		utils.RedirectError(w, r, "/register", "Some field not exist.", manager)
		return
	}
	if fields["password"] != fields["confirm_pass"] {
		utils.RedirectError(w, r, "/register", "The passwords don't match.", manager)
		return
	}
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}

	// Register new user.
	err = newAuth.RegisterUser(fields["username"], fields["password"])
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		utils.RedirectError(w, r, "/register", err.Error(), manager)
		return
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			utils.RedirectError(w, r, "/register", err.Error(), manager)
			return
		}
		_, err = db.SyncQ().Update("auth", []dbutils.DbEquals{{"name", fields["name"]}},
			[]dbutils.DbEquals{{"id", id}})
		if err != nil {
			utils.RedirectError(w, r, "/register", err.Error(), manager)
			return
		}
	} else {
		utils.RedirectError(w, r, "/register", fmt.Sprintf("Username %s not exist.", fields["username"]), manager)
		return
	}
	http.Redirect(w, r, "/sign-in", http.StatusFound)
}

func SignIn(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	manager.SetTemplatePath("src/templates/auth/signin.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": ""})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string)})
		manager.DelUserContext("error")
	}
	err := manager.RenderTemplate(w, r)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
}

func SignInPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	frm, err := utils.ParseForm(r)
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	fields, ok := utils.ConvertApplicationFormFields([]string{"username", "password"}, frm.GetApplicationForm())
	if !ok {
		utils.RedirectError(w, r, "/sign-in", "Some field not exist.", manager)
		return
	}
	// Connect to database.
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)

	// New auth.
	newAuth, err := auth.NewAuth(db)
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	err = newAuth.LoginUser(fields["username"], fields["password"])
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	// If the user exists set their ID in cookies.
	user, err := newAuth.UserExist(fields["username"])
	if err != nil {
		utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
		return
	}
	if user != nil {
		id, err := dbutils.ParseInt(user["id"])
		if err != nil {
			utils.RedirectError(w, r, "/sign-in", err.Error(), manager)
			return
		}
		cookie := &http.Cookie{
			Name:     "UID",
			Value:    strconv.Itoa(id),
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func ProfileView(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	_tempUID, _ := manager.GetUserContext("UID")
	UID, _ := strconv.Atoi(_tempUID.(string))

	id, _ := manager.GetSlugParams("id")
	err := conf.DatabaseI.Connect()
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	defer func(DatabaseI *database.Database) {
		err := DatabaseI.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(conf.DatabaseI)
	user, err := conf.DatabaseI.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", id}}, 1)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}

	// Render 404 if user not found
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("User not found"))
		if err != nil {
			utils.ServerError(w, err.Error())
			return
		}
		return
	}
	intId, _ := strconv.Atoi(id)
	userD := userData{
		Id:          intId,
		Name:        dbutils.ParseString(user[0]["name"]),
		Username:    dbutils.ParseString(user[0]["username"]),
		Avatar:      dbutils.ParseString(user[0]["avatar"]),
		Description: dbutils.ParseString(user[0]["description"]),
	}
	manager.SetTemplatePath("src/templates/profile.html")
	manager.SetContext(map[string]interface{}{"user": userD, "UID": UID})
	err = manager.RenderTemplate(w, r)
	if err != nil {
		panic(err)
	}
}

func ProfileEdit(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	if !editPermission(w, r, manager) {
		return
	}
	db := conf.DatabaseI
	err := db.Connect()
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(db)
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		utils.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id"))
		return
	}
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	user, err := db.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", uidInt}}, 1)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	userD := userData{
		Id:          uidInt,
		Name:        dbutils.ParseString(user[0]["name"]),
		Username:    dbutils.ParseString(user[0]["username"]),
		Avatar:      dbutils.ParseString(user[0]["avatar"]),
		Description: dbutils.ParseString(user[0]["description"]),
	}
	manager.SetTemplatePath("src/templates/profile_edit.html")
	myError, ok := manager.GetUserContext("error")
	manager.SetContext(map[string]interface{}{"error": "", "user": userD})
	if ok {
		manager.SetContext(map[string]interface{}{"error": myError.(string), "user": userD})
		manager.DelUserContext("error")
	}
	err = manager.RenderTemplate(w, r)
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
}

type editFormData struct {
	Name        string
	Description string
	DelAvatar   string
	AvatarPath  string
}

func ProfileEditPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	if !editPermission(w, r, manager) {
		return
	}
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		utils.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id"))
	}
	// Get form data
	frm, err := utils.ParseForm(r)
	if err != nil {
		utils.RedirectError(w, r, fmt.Sprintf("/profile/%s/edit", uid), err.Error(), manager)
		return
	}
	editForm := editFormData{
		Name:        frm.Value("name"),
		Description: frm.Value("description"),
		DelAvatar:   frm.Value("del_avatar"),
	}
	// Save avatar
	if editForm.DelAvatar == "" {
		file, fileHeader, err := frm.File("avatar")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			utils.ServerError(w, err.Error())
			return
		}
		var buildPath string
		if !errors.Is(err, http.ErrMissingFile) {
			if !SaveFile(w, file, fileHeader, "media/avatars/", &buildPath) {
				return
			}
		}
		editForm.AvatarPath = buildPath
	}
	// Update profile
	db := conf.DatabaseI
	err = db.Connect()
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
			return
		}
	}(db)
	updSlice := []dbutils.DbEquals{
		{"name", editForm.Name},
		{"description", editForm.Description},
	}
	// Delete avatar
	user, err := db.SyncQ().Select([]string{"*"}, "auth", []dbutils.DbEquals{{"id", uid}}, 1)
	dbAvatarPath := dbutils.ParseString(user[0]["avatar"])
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	if editForm.DelAvatar != "" {
		updSlice = append(updSlice, dbutils.DbEquals{Name: "avatar", Value: ""})
		if utils2.PathExist(dbAvatarPath) {
			err := os.Remove(dbAvatarPath)
			if err != nil {
				utils.ServerError(w, err.Error())
				return
			}
		}
	} else {
		// Delete old avatar
		if editForm.AvatarPath != "" {
			updSlice = append(updSlice, dbutils.DbEquals{Name: "avatar", Value: editForm.AvatarPath})
			if utils2.PathExist(dbAvatarPath) {
				err := os.Remove(dbAvatarPath)
				if err != nil {
					utils.ServerError(w, err.Error())
					return
				}
			}
		}
	}
	_, err = db.SyncQ().Update("auth", updSlice, []dbutils.DbEquals{{"id", uid}})
	if err != nil {
		utils.ServerError(w, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/prof/%s", uid), http.StatusFound)
}

func ProfileLogOutPost(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	cookie := &http.Cookie{
		Name:     "UID",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/sign-in", http.StatusFound)
}

func saveFileExist(pathToDir string, fileName string) string {
	outputFilepath := pathToDir + fileName
	if utils2.PathExist(pathToDir + fileName) {
		hash := sha256.Sum256([]byte(fileName))
		hashData := hex.EncodeToString(hash[:])
		ext := filepath.Ext(fileName)
		return saveFileExist(pathToDir, hashData+ext)
	}
	return outputFilepath
}

func SaveFile(w http.ResponseWriter, file multipart.File, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) bool {
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(file)

	fp := saveFileExist(pathToDir, fileHeader.Filename)
	*buildPath = fp
	dst, err := os.Create(fp)
	if err != nil {
		utils.ServerError(w, err.Error())
		return false
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			utils.ServerError(w, err.Error())
		}
	}(dst)
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.ServerError(w, err.Error())
		return false
	}
	return true
}

func editPermission(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) bool {
	uid, ok := manager.GetSlugParams("id")
	if !ok {
		utils.ServerError(w, fmt.Sprintf("Error when retrieving slug parameter %s.", "id"))
		return false
	}
	uidC, err := r.Cookie("UID")
	if err != nil {
		utils.ServerError(w, err.Error())
		return false
	}
	if uid != uidC.Value {
		utils.ServerForbidden(w)
		return false
	}
	return true
}

type userData struct {
	Id          int
	Name        string
	Username    string
	Avatar      string
	Description string
}
