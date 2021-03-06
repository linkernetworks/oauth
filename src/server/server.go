package server

import (
	"net/http"

	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/osinstorage"
	"github.com/linkernetworks/oauth/src/webservice"
)

type OAuthServer struct {
	config config.GlobalConfig
	http.Server
	oauthWebService *webservice.Service
}

func New(c config.GlobalConfig) *OAuthServer {

	s := &OAuthServer{
		config: c,
	}

	logger.Debugf("config: %#v", s.config)

	return s
}

func (s *OAuthServer) ListenAndServe() error {

	// TODO: use mongo
	storage := osinstorage.NewMemoryStorage(
		// TODO: input client data from outside
		osin.DefaultClient{
			Id:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://localhost",
		},
	)

	osin := osin.NewServer(&s.config.OsinConfig, storage)

	// TODO: use redis
	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)

	oauthWebService, err := webservice.New(
		osin,
		store,
	)
	if err != nil {
		logger.Fatalf("Create OAuth web service failed. err: [%v]", err)
	}

	oauthContainer := restful.NewContainer()
	oauthContainer.Add(oauthWebService.WebService())

	r := mux.NewRouter()
	r.PathPrefix("/oauth2/").Handler(oauthContainer)

	s.Handler = r

	return s.Server.ListenAndServe()
}
