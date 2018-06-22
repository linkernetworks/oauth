package app

import (
	"net/http"

	"github.com/linkernetworks/oauth/src/entity"
	"github.com/linkernetworks/oauth/src/mongo"
	"github.com/linkernetworks/oauth/src/validator"
	"github.com/RangelReale/osin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// OAuthAuthorize OAuthAuthorize handler
func OAuthAuthorize(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	// redirectUrl := "/signin?" + c.Request.URL.RawQuery
	oauthQueries := r.URL.RawQuery
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	locale := MustLocaleFunc(r)

	logrus.Info(r.URL.Query())

	userEmail := r.URL.Query()["user_email"][0]
	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(userEmail)
	if err != nil {
		validations["email"] = emailValidate
	}
	if validations.HasError() {
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	// check user access token in session, redirect to signin if does not exist
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userAccessToken := session.Values[userEmail]
	if userAccessToken == nil {
		logrus.Warnln("authorize endpoint, user access token is nil, redirect to login page")
		http.Redirect(w, r, "singin.html"+oauthQueries, http.StatusMovedPermanently)
		return
	}

	var user entity.User
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"email": userEmail},
	}
	err = mongoContext.QueryOne(sql, &user)
	if err != nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrQueryUser),
		})
	}
	if user.IsTokenExpired(appService.OAuthConfig.ExpiryDuration) {
		logrus.Warnln("user access token is expired, redirect to login.")
		http.Redirect(w, r, "signin"+oauthQueries, http.StatusMovedPermanently)
		return
	}

	resp := appService.OsinServer.NewResponse()
	defer resp.Close()

	// generate authorize requst using osin func.ar is nil when runinto error,
	// and error info is set in resp
	if ar := appService.OsinServer.HandleAuthorizeRequest(resp, r); ar != nil {
		ar.Authorized = true
		appService.OsinServer.FinishAuthorizeRequest(resp, r, ar)
	}

	// check error HandleAuthorizeRequest encounter error, ar is nil
	if resp.IsError && resp.InternalError != nil {
		logrus.Warnln("handle authorize request: ", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
}
