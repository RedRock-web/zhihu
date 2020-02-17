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
			"status":     42,
			"error_info": "不是提问！",
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

//获取问题信息
func (q *Question) GetQuestion() {
	tempTitle, _ := database.G_DB.Table.Find("question", "title", "question_id", q.Id)
	q.Title = string(tempTitle[0]["title"].([]uint8))

	tempUid, _ := database.G_DB.Table.Find("question", "uid", "question_id", q.Id)
	q.Uid = string(tempUid[0]["uid"].([]uint8))

	tempTime, _ := database.G_DB.Table.Find("question", "time", "question_id", q.Id)
	q.Time = string(tempTime[0]["time"].([]uint8))

	tempComplement, _ := database.G_DB.Table.Find("question", "complement", "question_id", q.Id)
	q.Complement = string(tempComplement[0]["complement"].([]uint8))
}

//获取问题的答案
func (q Question) GetAnswers() []map[string]interface{} {
	answers, err := database.G_DB.Table.HighFind("answer ", "uid,answer_id, time, content ", "question_id = "+q.Id)
	basic.CheckError(err, "获取问题的答案失败!")
	return answers
}

//获取答案数目
func (q Question) GetAnswersCount() int {
	counts, err := database.G_DB.Table.HignCount("answer ", "id", " question_id = "+q.Id)
	basic.CheckError(err, "回答评论子计数失败!")
	return counts
}

//搜索问题-通过正则表达式
func (q Question) Search(tableName string, targe string, limitInfo string) ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.HighFind(tableName, targe, "`title` regexp '"+limitInfo+"'")
	return data, err
}

//随机获取5条数据
func (q Question) GetByRand() ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.GetByRand("question", "question_id", "5")
	return data, err
}

//获取关注度最高的前n个问题
func (q Question) GetByFollowNum() ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.GetByOneself("select question_id,count(*) as count from question_follow group by question_id order by count desc limit 5")
	return data, err
}

