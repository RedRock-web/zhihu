package route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"zhihu/cmd/features"
)

func RoutePrepare(db *sql.DB) {
	r := gin.Default()

	LoginPage(db, r)
	Homepage(db, r)
	PersonalPage(db, r)
	IssueDetailsPage(db, r)

	r.Run()
}

func Homepage(db *sql.DB, r *gin.Engine) {
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
	r.GET("/logout", func(c *gin.Context) {
		features.Logout(db, c)
	})
}

//TODO:帐号重置密码
func LoginPage(db *sql.DB, r *gin.Engine) {
	r.POST("/sign_in", func(c *gin.Context) {
		features.RegisteOrLogin(db, c, "user_information")
	})
	r.GET("/account/password_reset", func(c *gin.Context) {

	})
	return
}

func PersonalPage(db *sql.DB, r *gin.Engine, username string) {
	r.PUT("/me", func(c *gin.Context) {
		features.AlterInformation(c, db, username)
	})
}

func IssueDetailsPage(db *sql.DB, r *gin.Engine, username string) {

}
