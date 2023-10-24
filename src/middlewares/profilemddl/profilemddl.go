package profilemddl

import (
	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"log"
	"net/http"
)

func AuthMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) {
	if r.URL.Path != "/register" && r.URL.Path != "/sign-in" && r.URL.Path != "/sign-in-post" && r.URL.Path != "/register-post" {
		uid, err := r.Cookie("UID")
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		// Connect to database.
		db := database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")
		err = db.Connect()
		if err != nil {
			log.Panicln(err)
		}
		defer func(db *database.Database) {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}(db)

		res, err := db.SyncQ().Select([]string{"username"}, "auth", []dbutils.DbEquals{{"id", uid.Value}}, 1)
		if err != nil {
			log.Panicln(err)
		}
		if res == nil {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
		}
	}
}