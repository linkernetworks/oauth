package app

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/linkernetworks/oauth/src/entity"
	"github.com/linkernetworks/oauth/src/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

var (
	verifyUserId = "59911e172f19"
	verifyCode   = "1234"
)

func newVerifyUser() entity.User {
	phoneNumber, ok := os.LookupEnv("TEST_SMS_PHONENUMBRR")
	if !ok {
		phoneNumber = "+8811111111111"
	}
	u := entity.User{
		ID:                    bson.ObjectId(verifyUserId),
		Email:                 "test@linker.com",
		Password:              "123456",
		Cellphone:             phoneNumber,
		VerificationCode:      verifyCode,
		AccessToken:           "oauth.test.access.token",
		AccessTokenExpiryTime: util.GetCurrentTimestamp() + 3600,
	}

	return u
}

func getUserVerifyCookie() *http.Cookie {
	resp, err := http.Get(apiBaseUrl + "/v1/user/verify/request")
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		logrus.Fatalf("cookies length expected to be 1")
	}

	return cookies[0]
}

func TestUserVerify(t *testing.T) {
	// EXECUTOR_NUMBER is a jenkins environment variable
	if os.Getenv("EXECUTOR_NUMBER") != "" {
		t.Skip("Fix this for concurrent build")
	}

	setupAPIServer(t)
	testLocale := testLocaleFunc(DefaultLocaleLang)
	user1 := newVerifyUser()
	upsertUser(user1)
	userVerifyAPI := apiBaseUrl + "/v1/user/verify"

	// test post with no cookie
	ret1 := PostForm(userVerifyAPI, nil)
	assert.NotNil(t, ret1)
	m1 := ret1.(map[string]interface{})
	assert.Equal(t, m1["message"], testLocale(ErrCookieNoUserId))
	assert.Equal(t, m1["error"], true)

	// set cookie
	cookie := getUserVerifyCookie()
	// test post without code
	ret2 := PostFromWithCookie(userVerifyAPI, nil, cookie)
	m2 := ret2.(map[string]interface{})
	assert.Equal(t, m2["message"], testLocale(ErrNoVerificationCode))
	assert.Equal(t, m2["error"], true)

	// test post with sms code, code not equal
	errRequestValues := url.Values{}
	errRequestValues.Set("code", "4321")
	ret3 := PostFromWithCookie(userVerifyAPI, errRequestValues, cookie)
	m3 := ret3.(map[string]interface{})
	assert.Equal(t, m3["message"], testLocale(ErrVerificationCodeNotEqual))
	assert.Equal(t, m3["error"], true)

	// test post with correct sms code
	correctRequestValues := url.Values{}
	correctRequestValues.Set("code", verifyCode)
	ret4 := PostFromWithCookie(userVerifyAPI, correctRequestValues, cookie)
	m4 := ret4.(map[string]interface{})
	assert.Equal(t, m4["message"], testLocale(UserVerifiedSuccess))
	assert.Equal(t, m4["error"], false)

	// test post with user has been verified
	requestValues := url.Values{}
	requestValues.Set("code", verifyCode)
	ret5 := PostFromWithCookie(userVerifyAPI, requestValues, cookie)
	m5 := ret5.(map[string]interface{})
	assert.Equal(t, m5["message"], testLocale(ErrUserVerified))
	assert.Equal(t, m5["error"], true)

	t.Logf("TestUserSignIn done, clear mongo db")
	clearDatabase()
}
