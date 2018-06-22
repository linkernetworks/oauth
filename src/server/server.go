package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/endpoint"
	"github.com/linkernetworks/oauth/src/log"
	"github.com/linkernetworks/oauth/src/service"
)

type OAuthServer struct {
	config config.GlobalConfig
	ep     *endpoint.Endpoint
}

func New(c ...config.GlobalConfig) service.ServiceI {

	s := &OAuthServer{
		config: config.DefaultConfig,
	}
	s.ep = endpoint.New(s)

	if len(c) > 0 {
		if err := mergo.Merge(&s.config, c[0], mergo.WithOverride); err != nil {
			logger.Fatalf("Merge config failed. err: %v", err)
		}
	}

	if err := log.Init(s.config.LoggerConfig); err != nil {
		logger.Warnf("Initialize logger failed. err: %v", err)
	}

	logger.Debugf("config: %#v", s.config)

	return s
}

func (s *OAuthServer) Start() error {
	s.startHTTP()
	return nil
}

func (s *OAuthServer) Shutdown(ctx context.Context) error {
	return nil
}

func (s *OAuthServer) startHTTP() {

	logger.Infoln("Starting HTTP server")

	go func() {

		enableV1, err := strconv.ParseBool(s.config.EnableV1)
		if err != nil {
			logger.Fatalf("Parse boll string failed. err: %v", err)
		}

		enableV2, err := strconv.ParseBool(s.config.EnableV2)
		if err != nil {
			logger.Fatalf("Parse boll string failed. err: %v", err)
		}

		r := mux.NewRouter()
		s.addCommonRoute(r)
		if enableV1 {
			s.addV1Route(r)
		}
		if enableV2 {
			s.addV2Route(r)
		}

		https, err := strconv.ParseBool(s.config.UseHTTPS)
		if err != nil {
			logger.Fatalf("Parse boll string failed. err: %v", err)
		}

		if https {
			bind := ":" + strconv.Itoa(s.config.HTTPSPort)
			logger.Infof("Starting HTTPS server on [%v]", bind)
			err := http.ListenAndServeTLS(bind, s.config.CertPublicKey, s.config.CertPrivateKey, r)
			if err != nil {
				logger.Fatalf("Start HTTPS server failed. err: %v", err)
			}
			logger.Infof("HTTPS server started on [%v]", bind)
		} else {
			bind := ":" + strconv.Itoa(s.config.HTTPPort)
			logger.Infof("Starting HTTP server on [%v]", bind)
			err := http.ListenAndServe(bind, r)
			if err != nil {
				logger.Fatalf("Start HTTP server failed. err: %v", err)
			}
			logger.Infof("HTTP server started on [%v]", bind)
		}
	}()
}

func (s *OAuthServer) addCommonRoute(r *mux.Router) {
	r.HandleFunc("/ping", s.ep.Ping)
}

func (s *OAuthServer) addV1Route(r *mux.Router) {

}

func (s *OAuthServer) addV2Route(r *mux.Router) {

}
