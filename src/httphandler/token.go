package httphandler

import (
	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin"
)

func Token(c *gin.Context) {

	server := c.MustGet("osinServer").(*osin.Server)

	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAccessRequest(resp, c.Request); ar != nil {
		ar.Authorized = true
		server.FinishAccessRequest(resp, c.Request, ar)
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}
