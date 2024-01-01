package conf

import "github.com/uwine4850/foozy/pkg/database"

func NewDb() *database.Database {
	return database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")
}

var LoadMessages = 10
