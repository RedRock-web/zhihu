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

type USER struct {
	Id           string
	Username     string
	Password     string
	Gender       string
	Nickname     string
	Introduction string
	Avatar       string
}

func HaveCookie(c *gin.Context, key string) bool {
	_, err := c.Cookie(key)
	return err == nil
}
