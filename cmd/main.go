package main

import (
	_ "github.com/go-sql-driver/mysql"
	"zhihu/cmd/database"
	"zhihu/cmd/route"
)

func main() {
	database.G_DB = database.Start()
	var r route.Route
	r.Start()
}
