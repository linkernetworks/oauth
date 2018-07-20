package webservice

import (
	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
	"github.com/gorilla/sessions"
)

type Service struct {
	osin  *osin.Server
	store sessions.Store
	web   *restful.WebService
}

func New(osin *osin.Server, store sessions.Store) (*Service, error) {

	s := &Service{
		osin:  osin,
		store: store,
		web:   &restful.WebService{},
	}

	s.web.Path("/oauth2").Produces(restful.MIME_JSON)
	s.web.Route(s.web.
		GET("/authorize").
		To(s.authorize),
	)
	s.web.Route(s.web.
		POST("/token").
		To(s.token),
	)

	return s, nil
}

func (s *Service) WebService() *restful.WebService {
	return s.web
}
