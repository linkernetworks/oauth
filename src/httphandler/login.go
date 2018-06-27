package httphandler

import (
	"fmt"
	"net/http"

	"github.com/linkernetworks/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("email is required"),
		})
		return
	}

	password := c.PostForm("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("password is required"),
		})
		return
	}

	// TODO: verify email/password
	user_id := email
	logger.Debugf("user [%v] login", user_id)

	session := sessions.Default(c)
	session.Set("user_id", user_id)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
