package question_details_page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//问题结构体
type Question struct {
	c          *gin.Context
	time       string //提问的时间
	questioner string //提问的人
	id         string //问题id
	title      string //问题标题
	complement string //问题补充
	comment    string //问题评论
	reply      string //问题的回答
}

func Start(c *gin.Context) {
	var q Question
	q.c = c
	q.id = c.Param("questionId")
	targe := c.Param("targe")
	switch targe {
	case "followers":
		{
			q.Follow()
		}
	case "comments":
		{

		}
	case "answers":
		{

		}

	}
}

//关注或取消关注问题
func (q Question) Follow() {
	if q.IsFollow() {
		err := database.G_DB.Table.Delete("question_follow", "uid = "+basic.G_UserID+" and question_id = "+q.id)
		basic.CheckError(err, "取消关注问题失败！")
		if err == nil {
			q.c.JSON(http.StatusOK, gin.H{
				"isFollow": "no",
			})
		}
	} else {
		err := database.G_DB.Table.Insert("question_follow", []string{"uid", "question_id"}, []string{basic.G_UserID, q.id})
		basic.CheckError(err, "关注问题失败！")
		if err == nil {
			q.c.JSON(http.StatusOK, gin.H{
				"isFollow": "yes",
			})
		} else {
			q.c.JSON(500,gin.H{
				"error":"出错！",
			})
		}
	}
}

//判断是否已经关注问题
func (q Question) IsFollow() bool {
	data, err := database.G_DB.Table.HighFind("question_follow", "id", "uid = "+basic.G_UserID+" and "+"question_id = "+q.id)
	basic.CheckError(err,"查询是否关注问题失败！")
	//fmt.Println("####")
	//fmt.Println(string(data[0]["id"].([]uint8)))
	return data[0]["id"].([]uint8) == nil
}

//写回答
func (q Question) Reply() {

}
