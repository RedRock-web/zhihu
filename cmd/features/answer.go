package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//Answer表示一个回答
type Answer struct {
	C          *gin.Context
	Id         string //回答id
	Time       string //回答时间
	Uid        string //回答的用户
	QuestionId string //回答的问题
	Content    string //回答的内容
}

//NewAnswer返回一个回答对象
func NewAnswer() *Answer {
	return &Answer{}
}

//写回答
func (a Answer) Post() error {
	err := database.G_DB.Table.Insert("answer", []string{"uid", "question_id", "answer_id", "time", "content"}, []string{a.Uid, a.QuestionId, a.Id, a.Time, a.Content})
	basic.CheckError(err, "回答问题失败！")
	return err
}

//判断是否已经写了回答
func (a Answer) HaveAnswer() bool {
	data, err := database.G_DB.Table.HighFind("answer", "id", "uid = "+a.Uid+" and question_id = "+a.QuestionId)
	basic.CheckError(err, "判断是否写了回答失败！")
	return data != nil
}

//删除回答
func (a Answer) Delete() error {
	return database.G_DB.Table.Delete("answer", "uid = "+a.Uid+" and question_id = "+a.QuestionId)
}

func (a Answer) GetId() string {
	data, err := database.G_DB.Table.HighFind("answer", "answer_id", "uid = "+a.Uid+" and question_id = "+a.QuestionId)
	basic.CheckError(err, "获取回答id失败！")
	return string(data[0]["answer_id"].([]uint8))
}

func (a Answer) GetTime() string {
	data, err := database.G_DB.Table.HighFind("answer", "time", "uid = "+a.Uid+" and question_id = "+a.QuestionId)
	basic.CheckError(err, "获取回答time失败！")
	return string(data[0]["time"].([]uint8))
}

func (a Answer) GetContent() string {
	data, err := database.G_DB.Table.HighFind("answer", "content", "uid = "+a.Uid+" and question_id = "+a.QuestionId)
	basic.CheckError(err, "获取回答content失败！")
	return string(data[0]["content"].([]uint8))
}

/*
//获取本问题，登录用户的回答id
func GetAnswerId(questionId string) string {
	data, err := database.G_DB.Table.HighFind("answer", "id", "uid = "+G_user.Info.Uid+" and question_id = "+questionId)
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
	a.id = GetAnswerId(G_question_id)
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

*/
