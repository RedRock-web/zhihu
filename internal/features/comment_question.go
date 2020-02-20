package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/internal/basic"
	"zhihu/internal/database"
)

// QuestionComment表示一个对问题的评论
type QuestionComment struct {
	C          *gin.Context
	QuestionId string //问题id
	Uid        string //用户id
	Id         string //评论的id
	Pid        string //对谁评论
	Time       string //评论时间
	Content    string //评论内容
}

//NewQuestionCommnet返回一个评论对象
func NewQuestionComment() *QuestionComment {
	return &QuestionComment{}
}

//发表问题评论
func (qc QuestionComment) Post() error {
	err := database.G_DB.Table.Insert("question_comment", []string{"comment_id", "uid", "question_id", "pid", "time", "content"}, []string{qc.Id, qc.Uid, qc.QuestionId, qc.Pid, qc.Time, qc.Content})
	basic.CheckError(err, "发表问题评论失败！")
	return err
}

//删除问题评论
func (qc QuestionComment) Delete() error {
	err := database.G_DB.Table.Delete("question_comment", "comment_id = "+qc.Id)
	basic.CheckError(err, "删除问题评论失败!")
	return err
}

//对问题评论计数
func (qc QuestionComment) GetCount() int {
	counts, err := database.G_DB.Table.Count("question_comment", "id")
	basic.CheckError(err, "问题评论计数失败!")
	return counts
}

//查看子评论数
func (qc QuestionComment) GetChildCount() int {
	counts, err := database.G_DB.Table.HignCount("question_comment ", "id", "pid = "+qc.Pid)
	basic.CheckError(err, "问题评论子计数失败!")
	return counts
}

//查看该问题的全部评论
func (qc QuestionComment) GetAllComment() []map[string]interface{} {
	comment, err := database.G_DB.Table.HighFind("question_comment ", "uid, comment_id, pid, time, content ", "question_id = "+qc.QuestionId)
	basic.CheckError(err, "查看问题评论失败!")
	return comment
}

//查看该问题评论的子评论
func (qc QuestionComment) GetChildComment() []map[string]interface{} {
	comment, err := database.G_DB.Table.HighFind("question_comment ", "uid, comment_id, pid, time, content ", "question_id = "+qc.QuestionId+" and pid = "+qc.Pid)
	basic.CheckError(err, "查看问题子评论失败!")
	return comment
}

//查看某评论的up数
func (qc QuestionComment) GetUpCounts() int {
	counts, err := database.G_DB.Table.HignCount("comment_vote ", "id", "comment_id = "+qc.Id+" and attitude = 1")
	basic.CheckError(err, "查看某评论的up数失败!")
	if err != nil {
		return 0
	} else {
		return counts
	}
}

//查看某评论的down数
func (qc QuestionComment) GetDownCounts() int {
	counts, err := database.G_DB.Table.HignCount("comment_vote ", "id", "comment_id = "+qc.Id+" and attitude = 0")
	basic.CheckError(err, "查看某评论的up数失败!")
	if err != nil {
		return 0
	} else {
		return counts
	}
}
