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

//json返回回答详情
func PostAnswers(c *gin.Context,data []map[string]interface{}) {
	var answer []gin.H

	for _, v := range data {
		//获取答案相关信息,将一个答案组合为一个gin.H
		answer = append(answer, gin.H{
			"uid":         string(v["uid"].([]uint8)),
			"question_id": string(v["question_id"].([]uint8)),
			"answer_id":   string(v["answer_id"].([]uint8)),
			"time":        string(v["time"].([]uint8)),
			"content":     string(v["content"].([]uint8)),
		})
	}

	//所有回答组合后,返回json
	c.JSON(200, gin.H{
		"status": 0,
		"data": gin.H{
			"answer": answer,
		},
	})
	c.Abort()
}
