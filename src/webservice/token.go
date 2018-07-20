package webservice

import (
	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
)

func (s *Service) token(req *restful.Request, resp *restful.Response) {

	oResp := s.osin.NewResponse()
	defer oResp.Close()

	if ar := s.osin.HandleAccessRequest(oResp, req.Request); ar != nil {

		// TODO: authentic & authorize

		ar.Authorized = true
		s.osin.FinishAccessRequest(oResp, req.Request, ar)
	}

	osin.OutputJSON(oResp, resp.ResponseWriter, req.Request)
}
