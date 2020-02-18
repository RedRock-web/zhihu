package main

import (
	_ "github.com/go-sql-driver/mysql"
	"zhihu/internal/database"
	"zhihu/internal/route"
)

func main() {
	database.Start()
	route.Start()
}
