package config

import (
	"github.com/linkernetworks/oauth/mongo"
	"github.com/linkernetworks/oauth/sms"
	"github.com/linkernetworks/redis"
)

// AppConfig
type AppConfig struct {
	OAuthConfig  OAuthConfig       `json:"oauth"`
	MongoConfig  mongo.MongoConfig `json:"mongo"`
	TwilioConfig sms.TwilioConfig  `json:"twilio"`
	Redis        redis.RedisConfig `json:"redis"`
}

// OAuthConfig
type OAuthConfig struct {
	ExpiryDuration int64            `json:"expiryDuration"` // how many seconds token will expired
	Host           string           `json:"host"`
	Port           string           `json:"port"`
	Encryption     EncryptionConfig `json:"encryption"`
}

type EncryptionConfig struct {
	Salt   string `json:"salt"`
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	Length int    `json:"length"`
}
