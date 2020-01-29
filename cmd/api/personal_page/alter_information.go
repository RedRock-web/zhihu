package personal_page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//更改个人信息
func AlterInfo(c *gin.Context) {
	targe := c.PostForm("targe")
	content := c.PostForm("content")
	if IsTargeCompliance(targe) {
		err := database.G_DB.Table.Alter("user", targe, content, "uid", basic.G_UserID)
		basic.CheckError(err, "更改个人信息失败！")
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":"修改成功！",
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "非法操作！",
			},
		})
	}
}
//判断targe是否合理
func IsTargeCompliance(include string) bool {
	return include == "gender" || include == "information" || include == "nickname" || include == "avatar"
}
