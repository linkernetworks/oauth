package httphandler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthorizedUserOrRedirect(c *gin.Context) {

	if checkSession(c) {
		c.Next()
		return
	} else {
		// TODO: append request URL to login page
		c.Redirect(http.StatusTemporaryRedirect, "/login?redirect_uri="+c.Request.RequestURI)
		c.Status(http.StatusTemporaryRedirect)
	}
}

func AuthorizedUser(c *gin.Context) {

	if checkSession(c) {
		c.Next()
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Not authorized",
		})
	}
}

func checkSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		return true
	}
	return false
}
