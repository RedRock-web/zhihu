package features

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//Answer表示一个回答
type Answer struct {
	C           *gin.Context
	id          string //回答id
	time        string //回答时间
	uid         string //回答的用户
	question_id string //回答的问题
	content     string //回答的内容
}

//NewAnswer返回一个回答对象
func NewAnswer() *Answer {
	return &Answer{}
}

//写回答接口
func PostAnswer(c *gin.Context) {
	a := NewAnswer()
	a.question_id = c.Param("questionId")
	a.time = basic.GetTimeNow()
	a.id = basic.GetAQuestionId()
	a.content = c.PostForm("content")
	a.uid = G_user.Info.Uid
	a.Post()
}

//写回答
func (a Answer) Post() {
	err := database.G_DB.Table.Insert("answer", []string{"uid", "quesiont_id", "answer_id", "time", "content"}, []string{a.uid, a.question_id, a.id, a.time, a.content})
	basic.CheckError(err, "回答问题失败！")
}

//判断是否已经写了回答接口
func HaveAnswer(questionId string) bool {
	a := NewAnswer()
	a.uid = G_user.Info.Uid
	a.question_id = questionId
	return a.HaveAnswer()
}

//判断是否已经写了回答
func (a Answer) HaveAnswer() bool {
	data, err := database.G_DB.Table.HighFind("answer", "id", "uid = "+a.uid+" and question_id = "+a.question_id)
	basic.CheckError(err, "判断是否写了回答失败！")
	return data[0]["id"].([]uint8) == nil
}

//获取本问题，登录用户的回答id
func GetAnswerId() string {
	data, err := database.G_DB.Table.HighFind("answer", "id", "uid = "+a.uid+" and question_id = "+a.question_id)
	basic.CheckError(err, "获取回答id失败！")
	return string(data[0]["id"].([]uint8))
}

//查看特定回答接口
func ViewAnswer(c *gin.Context) {
	a := NewAnswer()
	a.question_id = c.Param("question_id")
	a.id = c.Param("answer_id")
	if a.View() == nil {
		c.JSON(http.StatusOK, gin.H{
			"time":        a.time,
			"question_id": a.question_id,
			"answer_id":   a.id,
			"content":     a.content,
			"uid":         a.uid,
		})
	} else {
		c.JSON(500, gin.H{
			"error": "查询答案失败！",
		})
	}

}

//根据answer_id查看答案
func (a Answer) View() error {
	temp, err := database.G_DB.Table.HighFind("answer", "uid,time,content", "answer_id = "+a.id)
	basic.CheckError(err, "查询答案失败！")
	a.uid = string(temp[0]["uid"].([]uint8))
	a.time = string(temp[0]["time"].([]uint8))
	a.content = string(temp[0]["content"].([]uint8))
	return err
}

func DeleteAnswer(c *gin.Context) {
	a := NewAnswer()
	a.id = GetAnswerId()
	err := a.Delete()
	basic.CheckError(err, "删除回答失败！")
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除回答成功！",
		})
	} else {
		c.JSON(500, gin.H{
			"error": "删除回答失败！",
		})
	}
}

//删除回答
func (a Answer) Delete() error {
	return database.G_DB.Table.Delete("answer", "answer_id = "+a.id)
}
