package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"zhihu/cmd/database"
	"zhihu/cmd/registeOrLogin"
)

func main() {
    WebStart()
}

func WebStart()  {
	db := DatabasePrepare()
	RoutePrepare(db)
}

func DatabasePrepare() (db *sql.DB) {
	db1 := database.OpenDatabase("root", "root", "mysql")
	defer db1.Close()
	database.CreateDatabase(db1, "zhihu")
	db = database.OpenDatabase("root", "root", "zhihu")
	defer db.Close()
	database.CreateTable(db, "user", "username varchar(20), password varchar(20) ")

	return db
}

func RoutePrepare(db *sql.DB)  {
	r := gin.Default()
	r.POST("/registeOrLogin", func(c *gin.Context) {
		registeOrLogin.RegisteOrLogin(db, c, "user")
	})
	r.Run()
}

