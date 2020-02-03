package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/features"
)

//使用route结构体，把站点各个页面模块作为方法
type Route struct {
	auth *gin.RouterGroup
}

func (route Route) Start() {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg":  "sdf",
			"test": "sdfewtsf",
		})
	})

	route.LoginPage(r)
	route.auth = r.Group("")
	route.auth.Use(route.AuthRequired())
	{
		route.HomePage()
		route.PersonalPage()
		route.QuestionDetailsPage()
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

//主页
func (route Route) HomePage() {
	//主页
	route.auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "这里是主页！",
		})
	})
	//搜索
	route.auth.GET("/search", )
	//推荐
	route.auth.GET("/feed/topstory/recommend", )
	//关注
	route.auth.GET("/feed/topstory/follow_wonderful", )
	//热榜
	route.auth.GET("/feed/topstory/hot", )
	//提问
	route.auth.POST("/questions", features.Quiz)
	//注销
	route.auth.GET("/logout", features.Logout)
}

//登录注册页
func (route Route) LoginPage(r *gin.Engine) {
	r.GET("/sign_in", func(c *gin.Context) {
		//已登录,直接跳转主页
		c.String(200, "sdf")
		if features.IsLogin(c, "userID") {
			basic.Redirect(c, "/")
		} else { //没有登录，跳转到登录注册页
			c.JSON(200, gin.H{
				"msg": "成功来到登录注册页，现在可以登录，注册，找回密码了！",
			})
			r.POST("/sign_in", features.RegisteOrLogin)
			r.GET("/account/password_reset", features.PasswdReset)
		}
	})
}

//用户详情页
func (route Route) PersonalPage() {
	//编辑个人资料
	route.auth.GET("/edit", func(c *gin.Context) {
		route.auth.PUT("/me", features.Edit)
	})
	//查看用户状态
	route.auth.GET("/members/:nickname/:targe/:followingValus", )
	route.auth.POST("/chat", )
}

// 问题详情页
func (route Route) QuestionDetailsPage() {

	//进入问题详情页，获取问题信息
	route.auth.GET("/questions/:questionId/", func(c *gin.Context) {
		questionId := c.Param("questionId")

		//关注问题
		route.auth.POST("/questions/:questionId/followers", features.Follow)

		//取消关注问题
		route.auth.DELETE("/questions/:questionId/followers", features.CancelFollow)

		//查看问题评论
		route.auth.GET("/questions/:questionId/comments")

		//对问题发表评论
		route.auth.POST("/questions/:questionId/comments", features.PostQuestionComments)

		//删除评论
		route.auth.DELETE("/comments/:commentId", features.DeleteComment)

		//判断是否写了回答
		if features.HaveAnswer(questionId) {
			answerId := features.GetAnswerId()
			//查看自己的回答
			route.auth.GET("/questions/:questionId/"+answerId, features.ViewAnswer)
			//删除自己的回答
			route.auth.DELETE("/answers/" + answerId, features.DeleteAnswer)
		} else {
			//写回答
			route.auth.POST("/questions/:questionId/draft", features.PostAnswer)
		}

		//查看回答评论

		//对回答发表评论
		route.auth.POST("/questions/:questionId/answers/:answerId/comments", features.PostAnswerComments)
	})

}
