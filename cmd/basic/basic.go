package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println(errors.New(errorMsg))
	}
}

//user表格字段
type USER struct {
	Id           string
	Username     string
	Password     string
	Uid          string
	Gender       string
	Nickname     string
	Introduction string
	Avatar       string
	Question_id string
	Reply_id string
	Favorite_id string
	Followers_id string
	Concern_id string
	Article_id string
}

func HaveCookie(c *gin.Context, key string) bool {
	_, err := c.Cookie(key)
	return err == nil
}
