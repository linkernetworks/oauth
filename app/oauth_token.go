package app

import (
	"errors"

	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/entity"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/mongo"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/util"
	"github.com/RangelReale/osin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func OAuthToken(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	resp := appService.OsinServer.NewResponse()
	defer resp.Close()
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()

	if ar := appService.OsinServer.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			// check user table to validate user
			var user entity.User
			user.Email, user.Password = ar.Username, ar.Password
			user.Password = util.EncryptPassword(user.Password)
			sql := mongo.Selector{
				Collection: "user",
				Selector:   bson.M{"email": user.Email, "password": user.Password},
			}

			existed, err := mongoContext.CheckExist(sql, &user)
			if err != nil || !existed {
				resp.IsError = true
				resp.InternalError = errors.New("auth error")
			}
			ar.Authorized = true
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		case osin.ASSERTION:
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}
		appService.OsinServer.FinishAccessRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		logrus.Warnln("handle access request error: ", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
}
