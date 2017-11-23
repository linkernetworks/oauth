package mongo

import (
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// MongoClient
type MongoClient struct {
	session  *mgo.Session
	database string
}

// NewMongoClient
func NewMongoClient(mongoConfig MongoConfig) *MongoClient {
	addr := mongoConfig.Host + ":" + mongoConfig.Port
	dialInfo := mgo.DialInfo{
		Addrs:    []string{addr},
		Direct:   true,
		Database: mongoConfig.Database,
		Username: mongoConfig.User,
		Password: mongoConfig.Password,
		Timeout:  time.Duration(mongoConfig.TimeOut) * time.Second,
	}

	sess, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		logrus.Fatalf("Dial to mongo error: %v", err)
	}

	return &MongoClient{
		session:  sess,
		database: mongoConfig.Database,
	}
}

// NewContext is a way to copy new session for mgo
func (m *MongoClient) NewContext() *Context {
	sess := m.session.Copy()
	return &Context{
		client:   m,
		database: sess.DB(m.database),
		session:  sess,
	}
}
