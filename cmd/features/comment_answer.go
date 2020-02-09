package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

// AnswerComment表示一个对问题的评论
type AnswerComment struct {
	C        *gin.Context
	AnswerId string //回答id
	Uid      string //用户id
	Id       string //评论的id
	Pid      string //对谁评论
	Time     string //评论时间
	Content  string //评论内容
}

//NewAnswerCommnet返回一个评论对象
func NewAnswerComment() *AnswerComment {
	return &AnswerComment{}
}

//发表回答评论
func (ac AnswerComment) Post() error {
	err := database.G_DB.Table.Insert("answer_comment", []string{"comment_id", "uid", "answer_id", "pid", "time", "content"}, []string{ac.Id, ac.Uid, ac.AnswerId, ac.Pid, ac.Time, ac.Content})
	basic.CheckError(err, "发表问题评论失败！")
	return err
}

//删除回答评论
func (ac AnswerComment) Delete() error {
	err := database.G_DB.Table.Delete("answer_comment", "comment_id = "+ac.Id)
	basic.CheckError(err, "删除回答评论失败!")
	return err
}

//对回答评论计数
func (ac AnswerComment) GetCount() int {
	counts, err := database.G_DB.Table.Count("answer_comment", "id")
	basic.CheckError(err, "回答评论计数失败!")
	return counts
}

//查看该回答的全部评论
func (ac AnswerComment) GetAllComment() []map[string]interface{} {
	comment, err := database.G_DB.Table.HighFind("answer_comment", "uid, comment_id, pid, time, content ", "answer_id = "+ac.Id)
	basic.CheckError(err, "查看回答评论失败!")
	return comment
}
