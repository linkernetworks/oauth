package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	configPath := "../../config/test.json"
	// test get dev props
	appConfig := Read(configPath)

	assert.Equal(t, appConfig.MongoConfig.Database, "test_oauth")
	assert.Equal(t, appConfig.MongoConfig.Host, "localhost")
	assert.Equal(t, appConfig.OAuthConfig.ExpiryDuration, int64(3600))
}

func TestReadCurrentJSONConfig(t *testing.T) {
	if _, found := os.LookupEnv("OAUTH_CONFIG_PATH"); found {
		appConfig := ReadCurrent()
		assert.Equal(t, appConfig.MongoConfig.Database, "test_oauth")
		assert.Equal(t, appConfig.MongoConfig.Host, "localhost")
		assert.Equal(t, appConfig.OAuthConfig.ExpiryDuration, int64(3600))
	}
}
