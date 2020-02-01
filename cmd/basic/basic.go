package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

//检查error
func CheckError(err error, errorMsg string) {
	if err != nil {
		//fmt.Println(err)
		fmt.Println(errors.New(errorMsg))
	}
}

//根据cookie获取uid
func GetUid(c *gin.Context) string {
	uid, err := c.Cookie("userID")
	CheckError(err, "获取Uid失败！")
	return uid
}

//TODO:考虑各种算法获取用户唯一id
//获取当前时间-格式："2006-01-02 15:04:05"
func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//根据时间戳获取uid
func GetAUid() string {
	return strconv.FormatInt(time.Now().Unix()+98213523, 10)
}

//根据时间戳获取问题id
func GetAQuestionId() string {
	return strconv.FormatInt(time.Now().Unix()-1234567, 10)
}

//根据时间戳获取回答id
func GetAAnserId() string {
	return strconv.FormatInt(time.Now().Unix()+33911023, 10)
}

//根据时间戳获取评论id
func GetACommentId() string {
	return strconv.FormatInt(time.Now().Unix()-2951392, 10)
}

//重定向
func Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusMovedPermanently, url)
}
