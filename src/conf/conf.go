package conf

import "github.com/uwine4850/foozy/pkg/database"

var DatabaseI = database.NewDatabase("root", "1111", "mysql", "3406", "foozy_proj")

var LoadMessages = 5
