package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
	"zhihu/cmd/features"
)

func main() {
	WebStart()
}

func WebStart() {
	db := DatabasePrepare()
	RoutePrepare(db)
}

func DatabasePrepare() (db *sql.DB) {
	db1 := database.OpenDatabase("root", "root", "mysql")
	defer db1.Close()
	database.CreateDatabase(db1, "zhihu")
	db = database.OpenDatabase("root", "root", "zhihu")
	//defer db.Close()
	database.CreateTable(db, "user_information", "username varchar(20), password varchar(20) ,gender varchar(2), nickname varchar(2), introduction varchar(10), avatar varchar(10)")

	return db
}

func RoutePrepare(db *sql.DB) {
	var username string

	r := gin.Default()

	//主页
	r.GET("/", func(c *gin.Context) {

	})
	//搜索
	r.GET("search", func(c *gin.Context) {

	})
	//推荐
	r.GET("feed/topstory/recommend", func(c *gin.Context) {

	})
	//关注
	r.GET("feed/topstory/follow_wonderful", func(c *gin.Context) {

	})
	//热榜
	r.GET("feed/topstory/hot", func(c *gin.Context) {

	})
	//提问
	r.POST("/questions", func(c *gin.Context) {

	})

	r.POST("/sign_in", func(c *gin.Context) {
		username = features.RegisteOrLogin(db, c, "user")
	})

	r.GET("/logout", func(c *gin.Context) {
		features.Logout(db, c, username)
	})
	r.GET("/account/password_reset", func(c *gin.Context) {

	})
	r.PUT("/me", func(c *gin.Context) {
		AlterInformation(c, db, username)
	})

	r.Run()
}

func AlterInformation(c *gin.Context, db *sql.DB, username string) {
	targe := c.Query("include")
	content := c.Query(("content"))
	if TargeIsCompliance(targe) {
    	AlterContent(db, username, targe, content)
	}
}

func TargeIsCompliance(include string) bool {
	return include == "gender" || include == "imformation" || include == "nickname" || include == "avatar"
}
func AlterContent(db *sql.DB, username string, targe string, content string) {
	stmt, err := db.Prepare("update user_information set " +
		targe + " =" + content + " where username = " + username)
	basic.CheckError(err)
	stmt.Exec()
}


