package config

import (
	"bitbucket.org/linkernetworks/aurora/src/config"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/sms"
)

// AppConfig
type AppConfig struct {
	OAuthConfig  OAuthConfig        `json:"oauth"`
	MongoConfig  mongo.MongoConfig  `json:"mongo"`
	TwilioConfig sms.TwilioConfig   `json:"twilio"`
	Redis        config.RedisConfig `json:"redis"`
}

// OAuthConfig
type OAuthConfig struct {
	ExpiryDuration int64  `json:"expiryDuration"` // how many seconds token will expired
	Host           string `json:"host"`
	Port           string `json:"port"`
}
