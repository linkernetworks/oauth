package storage

import (
	"time"

	"github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/oauth/mongo"
	"github.com/RangelReale/osin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// methods fot osin mongostorage
const (
	CLIENT_COLLECTION    = "application"    // oauth2 client collection
	AUTHORIZE_COLLECTION = "authorizations" // oauth2 authorize collection
	ACCESS_COLLECTION    = "accesses"       // oauth2 accesses collection
)

type OsinStorage struct {
	session  *mgo.Session
	database string
}

func NewOsinStorage(mongoConfig mongo.MongoConfig) *OsinStorage {
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

	os := &OsinStorage{
		session:  sess,
		database: mongoConfig.Database,
	}

	return os
}

func (c *OsinStorage) Clone() osin.Storage {
	newsess := c.session.Copy()
	newc := &OsinStorage{
		session:  newsess,
		database: c.database,
	}

	return newc
}

func (c *OsinStorage) Close() {
}

// GetClient loads the client by id (client_id)
func (c *OsinStorage) GetClient(cliendID string) (osin.Client, error) {
	col := c.session.DB(c.database).C(CLIENT_COLLECTION)
	var client entity.Application
	err := col.Find(bson.M{"client_id": cliendID}).One(&client)

	return client, err
}

// SaveAuthorize saves authorize data.
func (c *OsinStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	col := c.session.DB(c.database).C(AUTHORIZE_COLLECTION)
	_, err := col.UpsertId(data.Code, data)

	return err
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (c *OsinStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	col := c.session.DB(c.database).C(AUTHORIZE_COLLECTION)
	var ad AuthorizeData
	err := col.Find(bson.M{"code": code}).One(&ad)

	return buildOsinAuthorizeData(&ad), err
}

// RemoveAuthorize revokes or deletes the authorization code.
func (c *OsinStorage) RemoveAuthorize(code string) error {
	col := c.session.DB(c.database).C(AUTHORIZE_COLLECTION)
	return col.RemoveId(code)
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (c *OsinStorage) SaveAccess(data *osin.AccessData) error {
	col := c.session.DB(c.database).C(ACCESS_COLLECTION)
	_, err := col.UpsertId(data.AccessToken, data)

	return err
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (c *OsinStorage) LoadAccess(accessToken string) (*osin.AccessData, error) {
	col := c.session.DB(c.database).C(ACCESS_COLLECTION)
	var ad AccessData
	err := col.Find(bson.M{"accesstoken": accessToken}).One(&ad)

	// build osin.AccessData
	accessData := &osin.AccessData{
		Client:        ad.Client,
		AuthorizeData: buildOsinAuthorizeData(ad.AuthorizeData),
		AccessToken:   ad.AccessToken,
		RefreshToken:  ad.RefreshToken,
		ExpiresIn:     ad.ExpiresIn,
		Scope:         ad.Scope,
		RedirectUri:   ad.RedirectUri,
		CreatedAt:     ad.CreatedAt,
		UserData:      ad.UserData,
	}
	accessData.AccessData = accessData

	return accessData, err
}

// RemoveAccess revokes or deletes an AccessData.
func (c *OsinStorage) RemoveAccess(token string) error {
	col := c.session.DB(c.database).C(ACCESS_COLLECTION)
	return col.RemoveId(token)
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (c *OsinStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	col := c.session.DB(c.database).C(ACCESS_COLLECTION)
	accessData := new(osin.AccessData)
	err := col.FindId(token).One(accessData)

	return accessData, err
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (c *OsinStorage) RemoveRefresh(token string) error {
	col := c.session.DB(c.database).C(ACCESS_COLLECTION)
	err := col.RemoveId(token)

	return err
}

func buildOsinAuthorizeData(ad *AuthorizeData) *osin.AuthorizeData {
	authorizeData := &osin.AuthorizeData{
		Client:              ad.Client,
		Code:                ad.Code,
		ExpiresIn:           ad.ExpiresIn,
		Scope:               ad.Scope,
		RedirectUri:         ad.RedirectUri,
		State:               ad.State,
		CreatedAt:           ad.CreatedAt,
		UserData:            ad.UserData,
		CodeChallenge:       ad.CodeChallenge,
		CodeChallengeMethod: ad.CodeChallengeMethod,
	}

	return authorizeData
}
