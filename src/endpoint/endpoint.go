package endpoint

import (
	"github.com/linkernetworks/oauth/src/service"
)

type Endpoint struct {
	s service.ServiceI
}

func New(s service.ServiceI) *Endpoint {
	return &Endpoint{
		s: s,
	}
}
