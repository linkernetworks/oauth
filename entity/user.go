package entity

import (
	"time"

	"bitbucket.org/linkernetworks/aurora/src/oauth/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID                    bson.ObjectId `bson:"_id" json:"id"`
	SerialNumber          string        `bson:"serial_number" json:"serial_number"`
	Email                 string        `bson:"email,omitempty" json:"email"`
	Password              string        `bson:"password,omitempty" json:"password,omitempty"`
	FirstName             string        `bson:"first_name" json:"first_name"`
	LastName              string        `bson:"last_name" json:"last_name"`
	CountryCode           string        `bson:"country_code" json:"country_code"`
	Cellphone             string        `bson:"cellphone" json:"cellphone"`
	Roles                 []string      `bson:"roles" json:"roles"`
	Verified              bool          `bson:"verified" json:"verified"`
	VerificationCode      string        `bson:"verification_code" json:"verification_code"`
	Jwt                   string        `bson:"jwt" json:"jwt"`
	AccessToken           string        `bson:"access_token" json:"access_token"`
	AccessTokenExpiryTime int64         `bson:"access_token_expiry_time" json:"access_token_expiry_time"`
	RefreshToken          string        `bson:"refresh_token" json:"refresh_token"`
	CreatedAt             int64         `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt             int64         `bson:"updated_at,omitempty" json:"updated_at"`
	LastLoggedInAt        int64         `bson:"last_loggedin_at,omitempty" json:"last_loggedin_at"`
	Revoked               bool          `bson:"revoked" json:"revoked"`
	JobPriority           float64       `bson:"job_priority" json:"job_priority"`
}

const USER_TOKEN_LENGTH = 24
const TOKEN_KEY = "LinkerNetworks.Inc"

// Implements VerificationProcessReceiver
func (u *User) SetVerificationCode(code string) {
	u.Verified = false
	u.VerificationCode = code
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) GetVerificationCode() string {
	return u.VerificationCode
}

func (u *User) SetVerified(val bool) bool {
	u.Verified = val
	return val
}

func (u *User) MatchVerificationCode(code string) bool {
	return u.VerificationCode == code
}

func (u *User) GetCellPhoneNumber() string {
	return u.Cellphone
}

func (u *User) GenerateToken(expiryDuration int64) {
	_, accessToken := util.GenPubPriKey(USER_TOKEN_LENGTH)
	_, refreshToken := util.GenPubPriKey(USER_TOKEN_LENGTH)

	u.AccessToken = accessToken
	u.RefreshToken = refreshToken
	u.AccessTokenExpiryTime = time.Now().Add(time.Duration(expiryDuration)).Unix()
	u.UpdatedAt = util.GetCurrentTimestamp()
}

// IsTokenExpired false user token not expired
// true user token expired
func (u *User) IsTokenExpired(expiryDuration int64) bool {
	if u.AccessToken == "" || u.AccessTokenExpiryTime == 0 {
		return true
	}

	// expiryDuration := props.GetInt64("oauth.user.expiry", 3600)
	if time.Now().Unix()-u.AccessTokenExpiryTime < expiryDuration {
		return false
	}

	return true
}

func (u *User) GenerateJwtToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenStr, err := token.SignedString([]byte(TOKEN_KEY))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (u *User) ParseJwtToken(tokenStr string) *jwt.Token {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_KEY), nil
	})
	if err != nil {
		logrus.Warnf("parse token error: %v, but ignore it.", err)
	}

	return token
}
