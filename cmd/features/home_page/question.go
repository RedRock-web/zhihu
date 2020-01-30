package home_page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//问题结构体
type Question struct {
	c          *gin.Context
	id         string
	time       string //提问的时间
	questioner string //提问的人
	title      string //问题标题
	complement string //问题补充
	comment    string //问题评论
	reply      string //问题的回答
}

func Start(c *gin.Context) {
	var q Question
	q.c = c
	q.Quiz()
}

//发起提问
func (q Question) Quiz() {
	q.title = q.c.PostForm("title")
	q.complement = q.c.PostForm("complement")
	q.time = basic.GetTimeNow()
	q.questioner = basic.G_UserID
	q.id = basic.GetAQuestionId()
	if q.IsQuestion() {
		err := database.G_DB.Table.Insert("question", []string{"uid", "time", "title", "complement", "question_id"}, []string{basic.G_UserID, q.time, q.title, q.complement, q.id})
		basic.CheckError(err, "提问失败！")
		if err == nil {
			q.c.JSON(http.StatusOK, gin.H{
				"question_id": q.id,
				"time":        q.time,
				"uid":         q.questioner,
				"title":       q.title,
				"complement":  q.complement,
			})
		} else {
			q.c.JSON(500, gin.H{
				"error": "提问失败！",
			})
		}
	} else {
		q.c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不是提问！",
		})
	}
}

//检查是否是提问，即是否以问号结尾
func (q Question) IsQuestion() bool {
	return q.title[len(q.title)-1:] == "?"
}
