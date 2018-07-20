package webservice

import (
	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
)

func (s *Service) authorize(req *restful.Request, resp *restful.Response) {

	oResp := s.osin.NewResponse()
	defer oResp.Close()

	if ar := s.osin.HandleAuthorizeRequest(oResp, req.Request); ar != nil {

		// TODO: authentic & authorize

		ar.Authorized = true
		s.osin.FinishAuthorizeRequest(oResp, req.Request, ar)
	}

	osin.OutputJSON(oResp, resp.ResponseWriter, req.Request)
}
