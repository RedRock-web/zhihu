package home_page

import "github.com/gin-gonic/gin"

func Quesion(c *gin.Context)  {
	title := c.PostForm("title")

}