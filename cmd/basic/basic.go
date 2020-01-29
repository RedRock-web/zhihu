package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)
var G_UserID string
var G_NickName string



func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println(errors.New(errorMsg))
	}
}

func GetUid(c *gin.Context) string {
	uid,err := c.Cookie("userID")
	CheckError(err, "获取Uid失败！")
	return uid
}