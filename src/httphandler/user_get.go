package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserGet(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"name":  "dummy user",
			"email": "eee@eee.eee",
		},
	})
}
