// +build integration

package app

import (
	"net/url"
	"os"
	"testing"

	"github.com/linkernetworks/oauth/src/util"
	"github.com/stretchr/testify/assert"
)

const TEST_SALT = "12345678"

func TestUserSignInI18n(t *testing.T) {

	// EXECUTOR_NUMBER is a jenkins environment variable
	if os.Getenv("EXECUTOR_NUMBER") != "" {
		t.Skip("Fix this for concurrent build")
	}

	setupAPIServer(t)
	userSininAPI := apiBaseUrl + "/v1/signin"
	userPost := url.Values{"email": {"test@linker.com"}, "password": {"123456"}}

	// test user not sign up and locale zh
	zhLocale := testLocaleFunc(ZHLocaleLang)
	zhHeader := make(map[string]string)
	zhHeader["Accept-Language"] = "zh"

	zhret1 := PostFormWithHeader(userSininAPI, userPost, zhHeader)
	zhm1 := zhret1.(map[string]interface{})
	assert.Equal(t, zhm1["message"], zhLocale(ErrUserNotExisted))
	assert.Equal(t, zhm1["error"], true)

	// add user without verified for test
	user := newTestUser()
	encrypted, err := util.EncryptPassword(user.Password, TEST_SALT)
	assert.NoError(t, err)
	user.Password = encrypted
	user.Verified = false
	upsertUser(user)

	// test user not verified
	zhret2 := PostFormWithHeader(userSininAPI, userPost, zhHeader)
	zhm2 := zhret2.(map[string]interface{})
	assert.Equal(t, zhm2["message"], zhLocale(ErrUserNotVerified))
	assert.Equal(t, zhm2["error"], true)

	// update user verified
	user.Verified = true
	upsertUser(user)

	// test user signin success
	zhret3 := PostFormWithHeader(userSininAPI, userPost, zhHeader)
	zhm3 := zhret3.(map[string]interface{})
	assert.Equal(t, zhm3["message"], zhLocale(UserLoginSuccess))

	// clear database
	clearDatabase()
	twuser := newTestUser()

	// test user not sign up and locale tw
	twLocale := testLocaleFunc(TWLocaleLang)
	twHeader := make(map[string]string)
	twHeader["Accept-Language"] = "zh-tw"

	twret1 := PostFormWithHeader(userSininAPI, userPost, twHeader)
	twm1 := twret1.(map[string]interface{})
	assert.Equal(t, twm1["message"], twLocale(ErrUserNotExisted))
	assert.Equal(t, twm1["error"], true)

	encrypted, err = util.EncryptPassword(twuser.Password, TEST_SALT)
	assert.NoError(t, err)
	twuser.Password = encrypted
	twuser.Verified = false
	upsertUser(twuser)

	// test user not verified
	twret2 := PostFormWithHeader(userSininAPI, userPost, twHeader)
	twm2 := twret2.(map[string]interface{})
	assert.Equal(t, twm2["message"], twLocale(ErrUserNotVerified))
	assert.Equal(t, twm2["error"], true)

	// update user verified
	twuser.Verified = true
	upsertUser(twuser)

	// test user signin success
	twret3 := PostFormWithHeader(userSininAPI, userPost, twHeader)
	twm3 := twret3.(map[string]interface{})
	assert.Equal(t, twm3["message"], twLocale(UserLoginSuccess))

	t.Logf("TestUserSignIn done, clear mongo db")
	clearDatabase()
}
