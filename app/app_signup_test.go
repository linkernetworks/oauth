package app

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppSignup(t *testing.T) {
	// server := setupAPIServer(t)
	setupAPIServer(t)
	testLocale := testLocaleFunc(DefaultLocaleLang)

	appSinupAPI := apiBaseUrl + "/v1/app/signup"
	appPost := url.Values{"name": {"linker"}}

	// test app first sign up
	ret1 := PostForm(appSinupAPI, appPost)
	m1 := ret1.(map[string]interface{})
	assert.Equal(t, m1["message"], testLocale(SaveAppSuccess))
	assert.Equal(t, m1["name"], "linker")

	// check app name existed
	ret2 := PostForm(appSinupAPI, appPost)
	m2 := ret2.(map[string]interface{})
	assert.Equal(t, m2["message"], testLocale(ErrAppNameExisted))
	assert.Equal(t, m2["error"], true)
}
