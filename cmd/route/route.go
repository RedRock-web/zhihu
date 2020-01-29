package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/api/login_page"
	"zhihu/cmd/api/personal_page"
)

//使用route结构体，把站点各个页面模块作为方法
type Route struct {
	auth *gin.RouterGroup
}

func (route Route) Start() {
	r := gin.Default()
	route.LoginPage(r)
	route.auth = r.Group("")
	route.auth.Use(route.AuthRequired())
	{
		route.HomePage()
		route.PersonalPage()
		route.IssueDetailsPage()
	}
	r.Run()
}


//使用cookie检测是否登录的中间件
func (route Route) AuthRequired() gin.HandlerFunc {
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

func (route Route) HomePage() {
	//主页
	route.auth.GET("/",)
	//搜索
	route.auth.GET("/search",)
	//推荐
	route.auth.GET("/feed/topstory/recommend", )
	//关注
	route.auth.GET("/feed/topstory/follow_wonderful",)
	//热榜
	route.auth.GET("/feed/topstory/hot",)
	//提问
	route.auth.POST("/questions", )
	route.auth.GET("/logout", login_page.Logout)
}

//TODO:帐号重置密码
func (route Route) LoginPage(r *gin.Engine) {
	r.POST("/sign_in", login_page.Start)
	r.GET("/account/password_reset",)
}

func (route Route) PersonalPage() {
	//更改个人信息
	route.auth.PUT("/me", personal_page.AlterInfo)
	//查看用户信息
	route.auth.GET("/members/:nickname/:targe/:followingValus",personal_page.ViewInfo)
	route.auth.POST("/chat", personal_page.Chat)
}

func (route Route) IssueDetailsPage() {}
