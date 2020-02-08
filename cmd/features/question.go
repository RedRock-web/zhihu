package features

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

// Question表示一个问题
type Question struct {
	C          *gin.Context
	Time       string //提问的时间
	Uid        string //提问的人
	Id         string //问题id
	Title      string //问题标题
	Complement string //问题补充
	Reply      string //问题的回答
}

//NewQuestion 返回一个问题对象
func NewQuestion() *Question {
	return &Question{}
}

//取消关注问题
func (q Question) CancelFollow() (err error) {
	err = database.G_DB.Table.Delete("question_follow", "uid = "+G_user.Info.Uid+" and question_id = "+q.Id)
	basic.CheckError(err, "取消关注问题失败！")
	return err
}

//关注问题
func (q Question) Follow() (err error) {
	err = database.G_DB.Table.Insert("question_follow", []string{"uid", "question_id"}, []string{G_user.Info.Uid, q.Id})
	basic.CheckError(err, "关注问题失败！")
	return err
}

// 判断是否已经关注问题
func (q Question) IsFollow() bool {
	data, err := database.G_DB.Table.HighFind("question_follow", "id", "uid = "+G_user.Info.Uid+" and "+"question_id = "+q.Id)
	basic.CheckError(err, "查询是否关注问题失败！")
	return data != nil
}

// 提问接口
func Quiz(c *gin.Context) {
	q := NewQuestion()
	q.C = c
	q.Quiz()
}

// 发起提问
func (q Question) Quiz() {
	q.Title = q.C.PostForm("title")
	q.Complement = q.C.PostForm("complement")
	q.Time = basic.GetTimeNow()
	q.Uid = G_user.Info.Uid
	q.Id = basic.GetAQuestionId()
	if q.IsQuestion() {
		err := database.G_DB.Table.Insert("question", []string{"uid", "time", "title", "complement", "question_id"}, []string{q.Uid, q.Time, q.Title, q.Complement, q.Id})
		basic.CheckError(err, "提问失败！")
		if err == nil {
			q.C.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data": gin.H{
					"question_id": q.Id,
					"time":        q.Time,
					"uid":         q.Uid,
					"title":       q.Title,
					"complement":  q.Complement,
				},
			})
		} else {
			q.C.JSON(500, gin.H{
				"error": "提问失败！",
			})
		}
	} else {
		q.C.JSON(http.StatusUnauthorized, gin.H{
			"error": "不是提问！",
		})
	}
}

// 检查是否是提问，即是否以问号结尾
func (q Question) IsQuestion() bool {
	return q.Title[len(q.Title)-1:] == "?"
}

//删除问题
func (q Question) Delete() {

}
