package httphandler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthorizedUserOrRedirect(c *gin.Context) {

	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		c.Next()
		return

	}

	// TODO: append request URL to login page

	c.Redirect(http.StatusTemporaryRedirect, "/login?redirect_uri="+c.Request.RequestURI)
	return
}
