package route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/features"
)

func RoutePrepare(db *sql.DB) {
	r := gin.Default()

	LoginPage(db, r)

	auth := r.Group("")
	auth.Use(AuthRequired())
	{
		Homepage(db, auth)
		PersonalPage(db, auth)
		IssueDetailsPage(db, auth)
	}

	r.Run()
}

//使用cookie检测是否登录的中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("userID")
		if cookie == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录！"})
			c.Abort()
		} else {
			c.Next()
			//c.JSON(200, "have cookie")
		}
	}
}

func Homepage(db *sql.DB, auth *gin.RouterGroup) {
	//主页
	auth.GET("/", func(c *gin.Context) {

	})
	//搜索
	auth.GET("/search", func(c *gin.Context) {

	})
	//推荐
	auth.GET("/feed/topstory/recommend", func(c *gin.Context) {

	})
	//关注
	auth.GET("/feed/topstory/follow_wonderful", func(c *gin.Context) {

	})
	//热榜
	auth.GET("/feed/topstory/hot", func(c *gin.Context) {

	})
	//提问
	auth.POST("/questions", func(c *gin.Context) {

	})
	auth.GET("/logout", func(c *gin.Context) {
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

func PersonalPage(db *sql.DB, auth *gin.RouterGroup) {
	auth.PUT("/me", func(c *gin.Context) {
		features.AlterInformation(c, db)
	})
}

func IssueDetailsPage(db *sql.DB, auth *gin.RouterGroup) {

}
