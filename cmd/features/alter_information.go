package features

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/database"
)

func AlterInformation(c *gin.Context, db *sql.DB, username string) {
	targe := c.Query("targe")
	content := c.Query(("content"))
	if TargeIsCompliance(targe) {
		database.DatabaseUpdate(db, "user_information", "username", username, targe, content)
		userInformation := database.DatabaseSearch(db, "user_information", "username", username)
		c.JSON(http.StatusOK, gin.H{
			"id":          userInformation.Id,
			"nickname":    userInformation.Nickname,
			"introdution": userInformation.Introduction,
			"gender":      userInformation.Gender,
			"avatar":      userInformation.Avatar,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "非法操作！",
			},
		})
	}
}

func TargeIsCompliance(include string) bool {
	return include == "gender" || include == "information" || include == "nickname" || include == "avatar"
}
