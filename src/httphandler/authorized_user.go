package httphandler

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/linkernetworks/logger"
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

	if checkSession(c) || checkOAuthToken(c) {
		c.Next()
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Not authorized",
		})
		c.Abort()
	}
}

func checkSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		return true
	}
	return false
}

func checkOAuthToken(c *gin.Context) bool {
	storage := c.MustGet("osinStorage").(osin.Storage)

	token := c.PostForm("token")
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		return false
	}

	data, err := storage.LoadAccess(token)
	switch err {
	case nil:
	default:
		logger.Warnf("Get OAuth access data failed. err: [%v]", err)
		fallthrough
	case osin.ErrNotFound:
		return false
	}

	return !data.IsExpired()
}
