package app

import (
	"bitbucket.org/linkernetworks/aurora/src/oauth/app/config"
	"bitbucket.org/linkernetworks/aurora/src/oauth/app/storage"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/sms"
	"github.com/linkernetworks/redis"
	"github.com/RangelReale/osin"
	"github.com/gin-contrib/sessions"
)

// ServiceProvider wrape dependencies of oauth api
type ServiceProvider struct {
	MongoClient  *mongo.MongoClient
	SmsClient    *sms.SMSClient
	OAuthConfig  *config.OAuthConfig
	OsinServer   *osin.Server
	SessionStore sessions.CookieStore
}

func NewServiceProvider(appConfig *config.AppConfig) *ServiceProvider {
	return NewServiceProviderFromConfig(appConfig)
}

func NewRedisService(appConfig *config.AppConfig) *redis.Service {
	return redis.New(&appConfig.Redis)
}
func NewServiceProviderFromConfig(appConfig *config.AppConfig) *ServiceProvider {
	as := &ServiceProvider{
		OAuthConfig:  &appConfig.OAuthConfig,
		MongoClient:  mongo.NewMongoClient(appConfig.MongoConfig),
		SessionStore: sessions.NewCookieStore([]byte("something-very-secret")),
	}

	sid := appConfig.TwilioConfig.Sid
	if sid != "" {
		as.SmsClient = sms.NewSMSClientFromConfig(appConfig.TwilioConfig)
	}

	// init osin.Server
	osinStorage := storage.NewOsinStorage(appConfig.MongoConfig)
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	as.OsinServer = osin.NewServer(sconfig, osinStorage)

	return as
}
