package entity

import (
	"testing"

	"github.com/linkernetworks/oauth/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func NewApplication() *Application {
	app_id := "4d88e15b60f486e428412dc9"
	app := &Application{
		ID:          bson.ObjectId(app_id),
		Name:        "Test001",
		Description: "Test001 for test",
		RedirectUri: "http://localhost:9090/callback",
		CreatedAt:   util.GetCurrentTimestamp(),
	}
	app.GenerateToken()

	return app
}

func TestGetId(t *testing.T) {
	app := NewApplication()
	id := app.GetId()
	assert.Equal(t, id, app.ClientId)
}

func TestGetSecret(t *testing.T) {
	app := NewApplication()
	secrect := app.GetSecret()
	assert.Equal(t, secrect, app.ClientSecret)
}

func TestGetRedirectUri(t *testing.T) {
	app := NewApplication()
	redirectUri := app.GetRedirectUri()
	assert.Equal(t, redirectUri, app.RedirectUri)
}

func TestGetUserData(t *testing.T) {
	app := NewApplication()
	userData := app.GetUserData()
	t.Skipf("user data is %v", userData)
}
