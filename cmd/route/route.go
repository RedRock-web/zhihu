package route

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/features"
)

//Engine 表示一个路由
type Engine struct {
	r *gin.Engine
}

//NewEngine 返回一个包含logger和recovery middleWare的路由
func NewEngine() *Engine {
	return &Engine{gin.Default()}
}

//main函数路由准备
func Start() {
	e := NewEngine()
	e.MiddleWare()
	e.Page()
	e.r.Run()
}

//设置全局路由middleWare
func (e Engine) MiddleWare() {
	e.r.Use(Cors())           //解决跨域问题
	e.r.NoRoute(NoResponse()) //解决404页面问题
}

//路由开启各项page feature
func (e Engine) Page() {
	e.LoginPage()
	e.PersonalPage()
}

//登录注册页feature
func (e Engine) LoginPage() {
	loginPage := e.r.Group("", LoginPageJudgeAuth())
	{
		loginPage.POST("/sign_in", features.RegisteOrLogin)
		loginPage.GET("/sign_in")
		loginPage.GET("/account/password_reset", features.PasswdReset)
	}
}

//用户详情页
func (e Engine) PersonalPage() {
	personalPage := e.r.Group("", AuthRequired())
	{
		//编辑个人资料
		personalPage.GET("/edit")
		personalPage.PUT("/me", features.Edit)
		personalPage.POST("/chat", )
	}
}

/*
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

// 问题详情页
func (route Route) QuestionDetailsPage() {

	//进入问题详情页，获取问题信息
	route.auth.GET("/questions/:questionId/", func(c *gin.Context) {
		features.G_question_id = c.Param("questionId")

		//关注问题
		route.auth.POST("/questions/:questionId/followers", features.Follow)

		//取消关注问题
		route.auth.DELETE("/questions/:questionId/followers", features.CancelFollow)

		//查看评论
		route.auth.GET("/comments/:commentId", )

		//查看子评论
		route.auth.GET("/comments/:commentId/child_comments", )

		//对问题发表评论
		route.auth.POST("/questions/:questionId/comments", features.PostQuestionComments)

		//删除评论
		route.auth.DELETE("/comments/:commentId", features.DeleteComment)

		//点赞或反对评论
		route.auth.POST("comments/:commentId/actions/like")
		route.auth.DELETE("comments/:commentId/actions/like")
		route.auth.POST("comments/:commentId/actions/dislike")
		route.auth.DELETE("comments/:commentId/actions/dislike")

		//判断是否写了回答
		if features.HaveAnswer(features.G_question_id) {
			answerId := features.GetAnswerId(features.G_question_id)
			//查看自己的回答
			route.auth.GET("/questions/"+features.G_question_id+"/"+answerId, features.ViewAnswer)
			//删除自己的回答
			route.auth.DELETE("/answers/"+answerId, features.DeleteAnswer)
		} else {
			//写回答
			route.auth.POST("/questions/:questionId/draft", features.PostAnswer)
		}

		//对回答发表评论
		route.auth.POST("/questions/:questionId/answers/:answerId/comments", features.PostAnswerComments)

		//对回答的点赞
		route.auth.POST("/answers/:answerId/voters")
	})
}
*/
