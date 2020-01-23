package personal_page

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/database"
)

//更改个人信息
func AlterInformation(c *gin.Context, db *sql.DB) {
	username, _ := c.Cookie("userID")
	targe := c.Query("targe")
	content := c.Query(("content"))
	if IsTargeCompliance(targe) {
		if nil == database.DatabaseUpdate(db, "user_information", "username", username, targe, content) {
			userInformation, err := database.DatabaseSearch(db, "user_information", "username", username)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"id":          userInformation.Id,
					"nickname":    userInformation.Nickname,
					"introdution": userInformation.Introduction,
					"gender":      userInformation.Gender,
					"avatar":      userInformation.Avatar,
				})
			}

		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "非法操作！",
			},
		})
	}
}

func IsTargeCompliance(include string) bool {
	return include == "gender" || include == "information" || include == "nickname" || include == "avatar"
}
