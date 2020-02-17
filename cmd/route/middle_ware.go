package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/features"
)

// 处理跨域请求,支持options访问的全局middleWare
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

//点击登录注册时,用cookie判断是否已经登录,如果已经登录,则重定向到主页,无法再次登录注册
func Unauthorized2LoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			basic.RediRect(c, "/")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

//点击登录后才能点的url,用cookie判断是否已经登录,否则重定向到登录注册页
func Authorized2Some() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			features.G_user.Info.Uid, _ = c.Cookie("userID")
			c.Next()
		} else {
			basic.RediRect(c, "/sign_in")
			c.Abort()
		}
	}
}

//404 响应middleWare
func NoResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": 50,
			"error_msg":  "page not exists",
		})
	}
}

//登录后,只要有新的请求,便刷新cookie的middleWare
func RefreshCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("userID")
		if err == nil {
			c.SetCookie("userID", cookie, 1000, "/", "127.0.0.1", false, true)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录！"})
		}
	}
}

//登录注册
func RegisteOrLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAccount()
		a.Username = c.PostForm("username")
		a.OriginalPasswd = c.PostForm("password")
		//密码加密
		a.Password = basic.Get32Md5(a.OriginalPasswd)
		a.C = c
		if a.IsRegiste("user") && a.Login() == nil {
			basic.RediRect(c, "/")
		} else if !a.IsRegiste("user") && a.Registe() == nil && a.Login() == nil {
			basic.RediRect(c, "/")
		}
	}
}

//注销帐号
func LogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := c.Cookie("userID")
		basic.CheckError(err, "注销失败！")
		c.SetCookie("userID", result, -1, "/", "127.0.0.1", false, true)
		c.JSON(http.StatusFound, gin.H{
			"message": "账号已注销！",
		})
	}
}

//写回答中间件
func ReplyAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.Time = basic.GetTimeNow()
		a.Id = basic.GetAQuestionId()
		a.Content = c.PostForm("content")
		a.QuestionId = c.Param("questionId")
		if !a.HaveAnswer() && basic.MethodIsOk(c, "POST") && a.Post() == nil {
			//
		}
		c.Abort()
	}
}

//删回答中间件
func DeleteAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.QuestionId = c.Param("questionId")
		if a.HaveAnswer() && basic.MethodIsOk(c, "DELETE") && a.Delete() == nil {
			//
		}
		c.Abort()
	}
}

//查看回答中间件
func ViewAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.QuestionId = c.Param("questionId")
		if a.HaveAnswer() && basic.MethodIsOk(c, "GET") {
			a.Id = a.GetId()
			a.Content = a.GetContent()
			a.Time = a.GetTime()
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data": gin.H{
					"uid":         a.Uid,
					"question_id": a.QuestionId,
					"answer_id":   a.Id,
					"time":        a.Time,
					"content":     a.Content,
				},
			})
		}
		c.Abort()
	}
}

//关注问题
func FollowQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := features.NewQuestion()
		q.Id = c.Param("questionId")
		if q.IsFollow() && basic.MethodIsOk(c, "DELETE") && q.CancelFollow() == nil {
			c.JSON(http.StatusOK, gin.H{
				"isFollow": "no",
			})
		} else if !q.IsFollow() && basic.MethodIsOk(c, "POST") && q.Follow() == nil {
			c.JSON(http.StatusOK, gin.H{
				"isFollow": "yes",
			})
		}
		c.Abort()
	}
}

//发表问题评论
func PostQuestionComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		qc := features.NewQuestionComment()
		qc.QuestionId = c.Param("questionId")
		qc.Time = basic.GetTimeNow()
		qc.Pid = c.DefaultPostForm("pid", "0")
		qc.Content = c.PostForm("content")
		qc.Id = basic.GetACommentId()
		qc.Uid = features.G_user.Info.Uid
		if basic.MethodIsOk(c, "POST") && qc.Post() == nil {
			//
		}
		c.Abort()
	}
}

//发表回答评论
func PostAnswerComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac := features.NewAnswerComment()
		ac.AnswerId = c.Param("answerId")
		ac.Time = basic.GetTimeNow()
		ac.Pid = c.DefaultPostForm("pid", "0")
		ac.Content = c.PostForm("content")
		ac.Id = basic.GetACommentId()
		ac.Uid = features.G_user.Info.Uid
		if basic.MethodIsOk(c, "POST") && ac.Post() == nil {
			//
		}
		c.Abort()
	}
}

//删除问题评论
func DeteleQuestionComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		qc := features.NewQuestionComment()
		qc.Id = c.Param("commentId")
		if qc.Delete() == nil {
			//
		}
		c.Abort()
	}
}

//删除回答评论
func DeteleAnswerComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac := features.NewAnswerComment()
		ac.Id = c.Param("commentId")
		if ac.Delete() == nil {
			//
		}
		c.Abort()
	}
}

//对评论表态
func VoteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		cv := features.NewCommentVote()
		cv.Time = basic.GetTimeNow()
		cv.Id = c.Param("commentId")
		cv.Uid = features.G_user.Info.Uid
		vote := c.PostForm("vote")
		if vote == "up" {
			cv.Attitude = "1"
		} else if vote == "down" {
			cv.Attitude = "0"
		}
		attitude := cv.GetAttitude()
		if attitude == "" { //原来的态度:不关心
			if vote == "up" { //原来是不关心,现在赞成就直接赞成
				cv.Agree()
			} else if vote == "down" { //原来是不关心,现在反对就直接反对
				cv.Against()
			}
		} else if attitude == "0" { //原来的态度:反对
			if vote == "up" { //原来是反对,现在要赞成,就先取消反对再赞成
				cv.CancelAgainst()
				cv.Agree()
			} else if vote == "down" { //原来的态度:反对
				cv.CancelAgainst() //原来是反对,现在要反对,就是双击取消反对,变成不关心
			}
		} else if attitude == "1" { //原来的态度:赞成
			if vote == "up" { //原来是赞成,现在要赞成,就是双击取消赞成,变成不关心
				cv.CancelAgree()
			} else if vote == "down" { //原来是赞成,现在要反对,就先取消赞成再反对
				cv.CancelAgree()
				cv.Against()
			}
		}
		c.Abort()
	}
}

//对回答表态
func VoteAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		av := features.NewAnswerVote()
		av.Time = basic.GetTimeNow()
		av.Id = c.Param("answerId")
		av.Uid = features.G_user.Info.Uid
		vote := c.PostForm("vote")
		if vote == "up" {
			av.Attitude = "1"
		} else if vote == "down" {
			av.Attitude = "0"
		}
		attitude := av.GetAttitude()
		if attitude == "" { //原来的态度:不关心
			if vote == "up" { //原来是不关心,现在赞成就直接赞成
				av.Agree()
			} else if vote == "down" { //原来是不关心,现在反对就直接反对
				av.Against()
			}
		} else if attitude == "0" { //原来的态度:反对
			if vote == "up" { //原来是反对,现在要赞成,就先取消反对再赞成
				av.CancelAgainst()
				av.Agree()
			} else if vote == "down" { //原来的态度:反对
				av.CancelAgainst() //原来是反对,现在要反对,就是双击取消反对,变成不关心
			}
		} else if attitude == "1" { //原来的态度:赞成
			if vote == "up" { //原来是赞成,现在要赞成,就是双击取消赞成,变成不关心
				av.CancelAgree()
			} else if vote == "down" { //原来是赞成,现在要反对,就先取消赞成再反对
				av.CancelAgree()
				av.Against()
			}
		}
		c.Abort()
	}
}

//查看问题评论
func ViewQuestionComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		qc := features.NewQuestionComment()   //获取新评论
		qc.QuestionId = c.Param("questionId") //记住该问题id
		counts := qc.GetCount()
		comments := qc.GetAllComment()

		var data []gin.H
		data = append(data, gin.H{"commet_counts": counts})
		for i := 0; i < counts; i++ {
			qc.Id = string(comments[i]["comment_id"].([]uint8))
			upCounts := qc.GetUpCounts()
			downCounts := qc.GetDownCounts()

			data = append(data, gin.H{
				"id":         i + 1,
				"uid":        string(comments[i]["uid"].([]uint8)),
				"comment_id": qc.Id,
				"pid":        string(comments[i]["pid"].([]uint8)),
				"content":    string(comments[i]["content"].([]uint8)),
				"time":       string(comments[i]["time"].([]uint8)),
				"vote_up":    upCounts,
				"vote_down":  downCounts,
			})
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   data,
		})
		c.Abort()
	}
}

//查看问题子评论
func ViewChildQuestionComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		qc := features.NewQuestionComment()   //获取新评论
		qc.QuestionId = c.Param("questionId") //记住该问题id
		qc.Pid = c.Param("commentId")
		counts := qc.GetChildCount()
		comments := qc.GetChildComment()
		var data []gin.H
		data = append(data, gin.H{"commet_counts": counts})
		for i := 0; i < counts; i++ {
			qc.Id = string(comments[i]["comment_id"].([]uint8))
			upCounts := qc.GetUpCounts()
			downCounts := qc.GetDownCounts()

			data = append(data, gin.H{
				"id":         i + 1,
				"uid":        string(comments[i]["uid"].([]uint8)),
				"comment_id": qc.Id,
				"pid":        string(comments[i]["pid"].([]uint8)),
				"content":    string(comments[i]["content"].([]uint8)),
				"time":       string(comments[i]["time"].([]uint8)),
				"vote_up":    upCounts,
				"vote_down":  downCounts,
			})
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   data,
		})
		c.Abort()
	}
}

//查看回答评论
func ViewAnswerComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac := features.NewAnswerComment() //获取新评论
		ac.AnswerId = c.Param("answerId") //记住该问题id
		counts := ac.GetCount()
		comments := ac.GetAllComment()

		var data []gin.H
		data = append(data, gin.H{"commet_counts": counts})
		for i := 0; i < counts; i++ {
			ac.Id = string(comments[i]["comment_id"].([]uint8))
			upCounts := ac.GetUpCounts()
			downCounts := ac.GetDownCounts()

			data = append(data, gin.H{
				"id":         i + 1,
				"uid":        string(comments[i]["uid"].([]uint8)),
				"comment_id": ac.Id,
				"pid":        string(comments[i]["pid"].([]uint8)),
				"content":    string(comments[i]["content"].([]uint8)),
				"time":       string(comments[i]["time"].([]uint8)),
				"vote_up":    upCounts,
				"vote_down":  downCounts,
			})
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   data,
		})
		c.Abort()
	}
}

//查看回答子评论
func ViewChildAnswerComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac := features.NewAnswerComment() //获取新评论
		ac.AnswerId = c.Param("answerId") //记住该问题id
		ac.Pid = c.Param("commentId")
		counts := ac.GetChildCount()
		comments := ac.GetChildComment()

		var data []gin.H
		data = append(data, gin.H{"commet_counts": counts})
		for i := 0; i < counts; i++ {
			ac.Id = string(comments[i]["comment_id"].([]uint8))
			upCounts := ac.GetUpCounts()
			downCounts := ac.GetDownCounts()

			data = append(data, gin.H{
				"id":         i + 1,
				"uid":        string(comments[i]["uid"].([]uint8)),
				"comment_id": ac.Id,
				"pid":        string(comments[i]["pid"].([]uint8)),
				"content":    string(comments[i]["content"].([]uint8)),
				"time":       string(comments[i]["time"].([]uint8)),
				"vote_up":    upCounts,
				"vote_down":  downCounts,
			})
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data":   data,
		})
		c.Abort()
	}
}

//获取问题详情
func GetQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := features.NewQuestion()
		q.Id = c.Param("questionId")
		q.GetQuestion()

		var question []gin.H
		question = append(question, gin.H{
			"author_uid":   q.Uid,
			"question_id":  q.Id,
			"created_time": q.Time,
			"title":        q.Title,
			"complement":   q.Complement,
		})

		var answer []gin.H
		tempAnswers := q.GetAnswers()
		answersCount := q.GetAnswersCount()
		for i := 0; i < answersCount; i++ {
			answer = append(answer, gin.H{
				"uid":         string(tempAnswers[i]["uid"].([]uint8)),
				"question_id": q.Id,
				"answer_id":   string(tempAnswers[i]["answer_id"].([]uint8)),
				"time":        string(tempAnswers[i]["time"].([]uint8)),
				"content":     string(tempAnswers[i]["content"].([]uint8)),
			})
		}
		c.JSON(200, gin.H{
			"status": 0,
			"data": gin.H{
				"question:": question,
				"answer":    answer,
			},
		})
		c.Abort()
	}
}

//提问
func Quiz() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := features.NewQuestion()
		q.C = c
		q.Quiz()
	}
}

//更改个人信息
func Edit() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := features.NewUser()
		user.Info.AlterTarge = c.PostForm("targe")
		user.Info.AlterContent = c.PostForm("content")
		user.Info.C = c
		if user.Info.IsTargeCompliance() {
			if user.Info.Alter() == nil && user.Info.View() == nil {
				c.JSON(http.StatusOK, gin.H{
					"status": 0,
					"data": gin.H{
						"uid":          user.Info.Uid,
						"nickname":     user.Info.Nickname,
						"gender":       user.Info.Gender,
						"avatar":       user.Info.Avatar,
						"introduction": user.Info.Introduction,
					},
				})
			}
		} else {
			user.Info.C.JSON(http.StatusUnauthorized, gin.H{
				"status":     41,
				"error_info": "非法操作!",
			})
		}
	}
}

//搜索feature
//TODO:按赞数排序
//TODO:搜索匹配回答
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		var answer []gin.H
		var question []gin.H

		qq := c.Query("q")
		q := features.NewQuestion()
		data, err := q.Search("question", "question_id", qq)
		basic.CheckError(err, "搜索问题失败!")

		if basic.MethodIsOk(c, "GET") && len(data) != 0 {
			for _, v := range data {
				//获取问题相关信息
				q.Id = string(v["question_id"].([]uint8))
				q.GetQuestion()
				//获取答案相关信息,将一个问题的答案组合为一个gin.H
				tempAnswers := q.GetAnswers()
				answersCount := q.GetAnswersCount()
				for i := 0; i < answersCount; i++ {
					answer = append(answer, gin.H{
						"uid":         string(tempAnswers[i]["uid"].([]uint8)),
						"question_id": q.Id,
						"answer_id":   string(tempAnswers[i]["answer_id"].([]uint8)),
						"time":        string(tempAnswers[i]["time"].([]uint8)),
						"content":     string(tempAnswers[i]["content"].([]uint8)),
					})
				}
				//将所有问题组合成一个gin.H
				question = append(question, gin.H{
					"author_uid":   q.Uid,
					"question_id":  q.Id,
					"created_time": q.Time,
					"title":        q.Title,
					"complement":   q.Complement,
					"answer":       answer,
				})
			}
			//所有问题组合后,返回json
			c.JSON(200, gin.H{
				"status": 0,
				"data": gin.H{
					"question:": question,
				},
			})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data":   "没有相关问题!",
			})
		}
	}
}

//主页,点击后随机获取5个问题
func HomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var answer []gin.H
		var question []gin.H

		q := features.NewQuestion()
		data, err := q.GetByRand()
		basic.CheckError(err, "搜索问题失败!")

		if basic.MethodIsOk(c, "GET") && len(data) != 0 {
			for _, v := range data {
				//获取问题相关信息
				q.Id = string(v["question_id"].([]uint8))
				q.GetQuestion()
				//获取答案相关信息,将一个问题的答案组合为一个gin.H
				tempAnswers := q.GetAnswers()
				answersCount := q.GetAnswersCount()
				for i := 0; i < answersCount; i++ {
					answer = append(answer, gin.H{
						"uid":         string(tempAnswers[i]["uid"].([]uint8)),
						"question_id": q.Id,
						"answer_id":   string(tempAnswers[i]["answer_id"].([]uint8)),
						"time":        string(tempAnswers[i]["time"].([]uint8)),
						"content":     string(tempAnswers[i]["content"].([]uint8)),
					})
				}
				//将所有问题组合成一个gin.H
				question = append(question, gin.H{
					"author_uid":   q.Uid,
					"question_id":  q.Id,
					"created_time": q.Time,
					"title":        q.Title,
					"complement":   q.Complement,
					"answer":       answer,
				})
			}
			//所有问题组合后,返回json
			c.JSON(200, gin.H{
				"status": 0,
				"data": gin.H{
					"question:": question,
				},
			})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data":   "没有问题推荐!",
			})
		}
	}
}
