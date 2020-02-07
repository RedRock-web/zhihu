package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
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
