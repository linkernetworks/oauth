package entity

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var testTokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAbGlua2VybmV0d29ya3MuY29tIiwiZXhwIjoxNTAxNTgyOTYyfQ.OwSVa-zGvKiIvJYoI3iG_IfFCD4-ASRXejcx60Uq750"
var testUser = User{
	Email:            "test@linkernetworks.com",
	Password:         "123456",
	VerificationCode: "oldcode",
	Verified:         true,
	Cellphone:        "+8615651006666",
}

func TestUpdateUserToken(t *testing.T) {
	u := testUser
	u.GenerateToken(3600)
	if u.AccessToken == "" || u.RefreshToken == "" {
		t.Errorf("build user error, AccessToken or RefreshToken should not empty")
	}
}

func TestSetVerificationCode(t *testing.T) {
	u := testUser
	u.SetVerificationCode("newcode")
	assert.Equal(t, u.GetVerificationCode(), "newcode")
	assert.Equal(t, u.Verified, false)
}

func TestSetVerified(t *testing.T) {
	u := testUser
	u.SetVerified(false)
	assert.Equal(t, u.Verified, false)
}

func TestMatchVerificationCode(t *testing.T) {
	u := testUser
	ret1 := u.MatchVerificationCode("oldcode")
	ret2 := u.MatchVerificationCode("newcode")
	assert.Equal(t, ret1, true)
	assert.Equal(t, ret2, false)
}

func TestGetCellPhoneNumber(t *testing.T) {
	u := testUser
	assert.Equal(t, u.GetCellPhoneNumber(), "+8615651006666")
}

func TestGetVerificationCode(t *testing.T) {
	u := User{
		Email:            "test@linkernetworks.com",
		Password:         "123456",
		VerificationCode: "oldcode",
	}
	assert.Equal(t, u.GetVerificationCode(), "oldcode")
}

func TestIsTokenExpired(t *testing.T) {
	u1 := User{
		Email:    "test@linkernetworks.com",
		Password: "123456",
	}
	u2 := User{
		Email:        "test@linkernetworks.com",
		Password:     "123456",
		AccessToken:  "aa",
		RefreshToken: "bb",
	}
	u3 := User{
		Email:                 "test@linkernetworks.com",
		Password:              "123456",
		AccessToken:           "aa",
		RefreshToken:          "bb",
		AccessTokenExpiryTime: time.Now().Unix() - 4600,
	}
	u4 := User{
		Email:                 "test@linkernetworks.com",
		Password:              "123456",
		AccessToken:           "aa",
		RefreshToken:          "bb",
		AccessTokenExpiryTime: time.Now().Unix(),
	}

	ret1 := u1.IsTokenExpired(3600) // true
	ret2 := u2.IsTokenExpired(3600) // true
	ret3 := u3.IsTokenExpired(3600) // true
	ret4 := u4.IsTokenExpired(3600) // false

	if !(ret1 && ret2 && ret3) {
		t.Error()
	}

	if ret4 {
		t.Error()
	}
}

func TestGenerateJwtToken(t *testing.T) {
	u := User{
		Email:    "test@linkernetworks.com",
		Password: "123456",
	}
	tokenStr, err := u.GenerateJwtToken()
	assert.NoError(t, err)
	if tokenStr == "" {
		t.Error()
	}
}

func TestParseJwtToken(t *testing.T) {
	u := User{
		Email:    "test@linkernetworks.com",
		Password: "123456",
	}
	token := u.ParseJwtToken(testTokenStr)
	assert.NotNil(t, token)

	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, claims["email"], "test@linkernetworks.com")
}
