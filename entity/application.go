package entity

import (
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/util"
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Name         string        `bson:"name" json:"name"`
	Description  string        `bson:"description" json:"description"`
	RedirectUri  string        `bson:"redirect_uri" json:"redirect_uri"`
	ClientId     string        `bson:"client_id" json:"client_id"`
	ClientSecret string        `bson:"client_secret" json:"client_secret"`
	CreatedAt    int64         `bson:"created_at" json:"created_at"`
	UpdatedAt    int64         `bson:"updated_at" json:"updated_at"`
	UserData     interface{}   `bson:"user_data" json:"user_data"`
}

func (a *Application) GenerateToken() {
	clientId, clientSecret := util.GenPubPriKey(12)
	a.ClientId = clientId
	a.ClientSecret = clientSecret
}

// GetId should return client_id
func (a Application) GetId() string {
	return a.ClientId
}

func (a Application) GetSecret() string {
	return a.ClientSecret
}

func (a Application) GetRedirectUri() string {
	return a.RedirectUri
}

func (a Application) GetUserData() interface{} {
	return a.UserData
}
