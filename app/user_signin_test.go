package app

import (
	"net/url"
	"testing"

	"bitbucket.org/linkernetworks/aurora/src/oauth/util"
	"github.com/stretchr/testify/assert"
)

func TestUserSignIn(t *testing.T) {
	setupAPIServer(t)
	testLocale := testLocaleFunc(DefaultLocaleLang)
	userSininAPI := apiBaseUrl + "/v1/signin"
	userPost := url.Values{"email": {"test@linker.com"}, "password": {"123456"}}

	// test user not sign up
	ret1 := PostForm(userSininAPI, userPost)
	m1 := ret1.(map[string]interface{})
	assert.Equal(t, m1["message"], testLocale(ErrUserNotExisted))
	assert.Equal(t, m1["error"], true)

	// add user without verified for test
	user := newTestUser()
	encrypted, err := util.EncryptPassword(user.Password)
	assert.NoError(t, err)
	user.Password = encrypted
	user.Verified = false
	upsertUser(user)

	// test user not verified
	ret2 := PostForm(userSininAPI, userPost)
	m2 := ret2.(map[string]interface{})
	assert.Equal(t, m2["message"], testLocale(ErrUserNotVerified))
	assert.Equal(t, m2["error"], true)

	// update user verified
	user.Verified = true
	upsertUser(user)

	// test user signin success
	ret3 := PostForm(userSininAPI, userPost)
	m3 := ret3.(map[string]interface{})
	assert.Equal(t, m3["message"], testLocale(UserLoginSuccess))

	t.Logf("TestUserSignIn done, clear mongo db")
	clearDatabase()
}
