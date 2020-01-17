package basic

import "github.com/gin-gonic/gin"

func CheckError(err error)  {
	if err != nil {
		panic(err)
	}
}








func HaveCookie(c *gin.Context, username string) bool {
	_, err := c.Cookie(username)
	return err == nil
}

