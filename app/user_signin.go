package app

import (
	"net/http"
	"net/url"

	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/entity"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/mongo"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/util"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// UserSignIn handle user signin logic
func UserSignIn(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	locale := MustLocaleFunc(r)

	var user entity.User
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(user.Email)
	if err != nil {
		validations["email"] = emailValidate
	}
	passworkValidate, err := validator.ValidatePassword(user.Password)
	if err != nil {
		validations["password"] = passworkValidate
	}

	if validations.HasError() {
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	// check if user existed in db
	user.Password = util.EncryptPassword(user.Password)
	sql := mongo.Selector{
		Collection: "user",
		Selector: bson.M{
			"email":    user.Email,
			"password": user.Password},
	}
	existed, err := mongoContext.CheckExist(sql, &user)
	if err != nil {
		logrus.Warnf("user sign in error: %v\n", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrCheckUserName),
		})
		return
	}
	if !existed {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserNotExisted),
		})
		return
	}

	if !user.Verified {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserNotVerified),
		})
		return
	}

	// check user access token, secrect token and expiry time
	expiryDuration := appService.OAuthConfig.ExpiryDuration
	if user.IsTokenExpired(expiryDuration) {
		logrus.Infof("user info expired, access_token: %s, refresh_token: %s, expired_ts: %d\n",
			user.AccessToken, user.RefreshToken, user.AccessTokenExpiryTime)
		user.GenerateToken(expiryDuration)
		logrus.Infof("new user info, access_token: %s, refresh_token: %s, expired_ts: %d\n",
			user.AccessToken, user.RefreshToken, user.AccessTokenExpiryTime)
		sql := mongo.Selector{
			Collection: "user",
			Selector:   bson.M{"_id": user.ID},
		}
		err := mongoContext.UpdateOne(sql, &user)
		if err != nil {
			logrus.Warnf("update user info error: %v", err)
			toJson(w, FormActionResponse{
				Error:   true,
				Message: locale(ErrUserUpdate),
			})
			return
		}
	}

	// set cookie and session before redirect to oauth server
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values[user.Email] = user.AccessToken
	session.Save(r, w)

	// judge if return to destination
	signinQuery := r.URL.RawQuery
	unescapeQuery, err := url.QueryUnescape(signinQuery)
	if err != nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrParseQuery),
		})
		return
	}
	values, err := url.ParseQuery(unescapeQuery)
	if err != nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrParseQuery),
		})
		return
	}
	if values.Get("response_type") != "" {
		redirecrURL := "/services/oauth/authorize?" + unescapeQuery
		http.Redirect(w, r, redirecrURL, http.StatusMovedPermanently)
		return
	}

	// user info all valide, return
	toJson(w, gin.H{
		"message":       locale(UserLoginSuccess),
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
		"expiry":        expiryDuration,
	})
}
