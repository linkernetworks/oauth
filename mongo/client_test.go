package mongo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newMongoConfig() MongoConfig {
	configPath := os.Getenv("OAUTH_CONFIG_PATH")
	if configPath == "" {
		configPath = "./test.json"
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("read config file %s error: %v\n", configPath, err)
	}

	var mongoConfig MongoConfig
	err = json.Unmarshal(content, &mongoConfig)
	if err != nil {
		logrus.Fatalf("parse mongo config error: %v\n", err)
	}

	return mongoConfig
}

func TestNewMongoClient(t *testing.T) {
	mongoConfig := newMongoConfig()
	mongoClient := NewMongoClient(mongoConfig)

	assert.NotNil(t, mongoClient)
	assert.NotNil(t, mongoClient.session)
	assert.Equal(t, mongoClient.database, mongoConfig.Database)
}
