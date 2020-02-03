package features

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

// Comment表示一个评论
type Comment struct {
	C       *gin.Context
	targeId string //targeId为questionId或answerId
	uid     string //用户id
	id      string //评论的id
	pid     string //对谁评论
	time    string //评论时间
	content string //评论内容
}

//NewCommnet返回一个评论对象
func NewComment() *Comment {
	return &Comment{}
}

//发表问题评论接口
func PostQuestionComments(c *gin.Context) {
	comment := NewComment()
	comment.targeId = c.Param("questionId")
	comment.time = basic.GetTimeNow()
	comment.pid = c.PostForm("pid")
	comment.content = c.PostForm("content")
	comment.id = basic.GetACommentId()
	comment.uid = G_user.Info.Uid
	comment.Post("question")
}

//发表回答评论接口
func PostAnswerComments(c *gin.Context) {
	comment := NewComment()
	comment.targeId = c.Param("answerId")
	comment.time = basic.GetTimeNow()
	comment.pid = c.PostForm("pid")
	comment.content = c.PostForm("content")
	comment.id = basic.GetACommentId()
	comment.uid = G_user.Info.Uid
	comment.Post("question")
}

//发表评论
func (c Comment) Post(targe string) {
	if targe == "question" {
		err := database.G_DB.Table.Insert("question_comment", []string{"comment_id", "uid", "question_id", "pid", "time", "content"}, []string{c.id, c.uid, c.targeId, c.pid, c.time, c.content})
		basic.CheckError(err, "发表问题评论失败！")
	} else if targe == "answer" {
		err := database.G_DB.Table.Insert("answer_comment", []string{"comment_id", "uid", "answer_id", "pid", "time", "content"}, []string{c.id, c.uid, c.id, c.targeId, c.time, c.content})
		basic.CheckError(err, "发表回答评论失败！")
	}
}

//删除问题评论
func (c Comment) DeleteQuestion() error {
	err := database.G_DB.Table.Delete("question_comment", "comment_id = "+c.id)
	basic.CheckError(err, "删除问题评论失败！")
	return err
}

//删除回答评论
func (c Comment) DeleteAnswer() error {
	err := database.G_DB.Table.Delete("answer_comment", "comment_id = "+c.id)
	basic.CheckError(err, "删除问题评论失败！")
	return err
}

//删除评论
func (c Comment) Delete() {
	if c.IsQuestionComment() {
		if c.DeleteQuestion() == nil {
			c.C.JSON(http.StatusOK, gin.H{
				"msg": "删除成功！",
			})
		} else {
			c.C.JSON(500, gin.H{
				"error": "删除失败！",
			})
		}
	} else if c.IsAnswerComment() {
		if c.DeleteAnswer() == nil {
			c.C.JSON(http.StatusOK, gin.H{
				"msg": "删除成功！",
			})
		} else {
			c.C.JSON(500, gin.H{
				"error": "删除失败！",
			})
		}
	} else {
		c.C.JSON(404, gin.H{
			"error": "评论不存在",
		})
	}
}

//删除评论接口
func DeleteComment(c *gin.Context) {
	comment := NewComment()
	comment.C = c
	comment.id = c.Param("commentId")
	comment.Delete()
}

//判断评论是否是问题的评论
func (c Comment) IsQuestionComment() bool {
	data, err := database.G_DB.Table.Find("question_comment", "id", "comment_id", c.id)
	basic.CheckError(err, "判断评论是否是问题的评论失败！")
	return data[0]["id"].([]uint8) == nil
}

//判断评论是否是回答的评论
func (c Comment) IsAnswerComment() bool {
	data, err := database.G_DB.Table.Find("answer_comment", "id", "comment_id", c.id)
	basic.CheckError(err, "判断评论是否是回答的评论失败！")
	return data[0]["id"].([]uint8) == nil
}
