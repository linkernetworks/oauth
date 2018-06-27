package config

import (
	"github.com/RangelReale/osin"
	"github.com/linkernetworks/logger"
)

type GlobalConfig struct {
	EnableV1       string
	EnableV2       string
	UseHTTPS       string
	HTTPPort       int
	HTTPSPort      int
	CertPublicKey  string
	CertPrivateKey string
	LoggerConfig   logger.LoggerConfig
	OsinConfig     osin.ServerConfig
}

var DefaultConfig GlobalConfig = GlobalConfig{
	EnableV1:       "false",
	EnableV2:       "true",
	UseHTTPS:       "true",
	HTTPPort:       8080,
	HTTPSPort:      8443,
	CertPublicKey:  "./tls-key/server.crt",
	CertPrivateKey: "./tls-key/server.key",
	LoggerConfig: logger.LoggerConfig{
		Dir:           "./logs",
		Level:         "info",
		MaxAge:        "720h",
		SuffixPattern: ".%Y%m%d",
		LinkName:      "oauth_log",
	},
	OsinConfig: osin.ServerConfig{
		AuthorizationExpiration: 250,
		AccessExpiration:        3600,
		TokenType:               "Bearer",
		AllowedAuthorizeTypes: []osin.AuthorizeRequestType{
			"code",
		},
		AllowedAccessTypes: []osin.AccessRequestType{
			"authorization_code",
		},
		ErrorStatusCode:             400,
		AllowClientSecretInParams:   false,
		AllowGetAccessRequest:       false,
		RequirePKCEForPublicClients: false,
		RedirectUriSeparator:        "",
		RetainTokenAfterRefresh:     false,
	},
}

var DevelopConfig GlobalConfig = GlobalConfig{
	UseHTTPS: "false",
	LoggerConfig: logger.LoggerConfig{
		Level: "debug",
	},
}
