package main

import (
	_ "github.com/go-sql-driver/mysql"
	"zhihu/cmd/database"
	"zhihu/cmd/route"
)

func main() {
	db := database.DatabasePrepare()
	route.RoutePrepare(db)
}
