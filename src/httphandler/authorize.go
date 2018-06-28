package httphandler

import (
	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
	server := c.MustGet("osinServer").(*osin.Server)

	resp := server.NewResponse()
	defer resp.Close()

	if ar := server.HandleAuthorizeRequest(resp, c.Request); ar != nil {

		// TODO: authentic & authorize

		ar.Authorized = true
		server.FinishAuthorizeRequest(resp, c.Request, ar)
	}

	osin.OutputJSON(resp, c.Writer, c.Request)
}
