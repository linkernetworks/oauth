package server

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/RangelReale/osin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/imdario/mergo"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/httphandler"
	"github.com/linkernetworks/oauth/src/log"
	"github.com/linkernetworks/oauth/src/osinstorage"
)

type OAuthServer struct {
	config      config.GlobalConfig
	router      *gin.Engine
	httpServer  *http.Server
	osinStorage *osinstorage.MemoryStorage
	osinServer  *osin.Server
}

func New(c ...config.GlobalConfig) *OAuthServer {

	s := &OAuthServer{
		config: config.DefaultConfig,
	}

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

	s.router = gin.Default()

	s.osinStorage = osinstorage.NewMemoryStorage(
		// TODO: input client data from outside
		osin.DefaultClient{
			Id:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://localhost",
		},
	)

	s.osinServer = osin.NewServer(&s.config.OsinConfig, s.osinStorage)
	s.router.Use(func(c *gin.Context) {
		c.Set("osinServer", s.osinServer)
		c.Next()
	})

	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)
	s.router.Use(sessions.Sessions("session", store))

	s.startHTTP()

	return nil
}

func (s *OAuthServer) Shutdown(ctx context.Context) error {

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logger.Warnf("Shutdown HTTP server failed. err: [%v]", err)
		}
		wg.Done()
	}()

	wait := make(chan int)
	go func() {
		wg.Wait()
		wait <- 1
	}()

	select {
	case <-ctx.Done():
	case <-wait:
	}

	return ctx.Err()
}

func (s *OAuthServer) startHTTP() error {

	s.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	s.router.GET("/login", httphandler.LoginPage)

	api := s.router.Group("/api")
	{
		api.POST("/login", httphandler.Login)

		oauthv2 := api.Group("/oauth2")
		{
			oauthv2.Use(httphandler.CheckAuthorizedUser).GET("/authorize", httphandler.Authorize)
			oauthv2.POST("/token", httphandler.Token)
		}
	}

	https, err := strconv.ParseBool(s.config.UseHTTPS)
	if err != nil {
		logger.Fatalf("Parse boll string failed. err: %v", err)
	}

	if https {
		bind := ":" + strconv.Itoa(s.config.HTTPSPort)
		logger.Infof("Starting HTTPS server on [%v]", bind)

		s.httpServer = &http.Server{
			Addr:    bind,
			Handler: s.router,
		}

		go func() {
			err := s.httpServer.ListenAndServeTLS(s.config.CertPublicKey, s.config.CertPrivateKey)
			if err != nil && err != http.ErrServerClosed {
				logger.Fatalf("Start HTTPS server failed. err: %v", err)
			} else {
				logger.Infof("HTTPS server closed.")
			}
		}()

	} else {
		bind := ":" + strconv.Itoa(s.config.HTTPPort)
		logger.Infof("Starting HTTP server on [%v]", bind)

		s.httpServer = &http.Server{
			Addr:    bind,
			Handler: s.router,
		}

		go func() {
			err := s.httpServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logger.Fatalf("Start HTTP server failed. err: %v", err)
			} else {
				logger.Infof("HTTP server closed.")
			}
		}()
	}

	return nil
}
