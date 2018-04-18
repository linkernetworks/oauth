package app

import (
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/verification"
	"bitbucket.org/linkernetworks/aurora/src/service/provider"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// UserVerify verify user by sms verification code.
func UserVerify(w http.ResponseWriter, r *http.Request, appService *provider.Container) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	locale := MustLocaleFunc(r)

	cookieUserID := session.Values["current_user_id"]
	if cookieUserID == nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrCookieNoUserId),
		})
		return
	}
	var cookieUserIDStr = cookieUserID.(string)
	var user entity.User
	user.ID = bson.ObjectId(cookieUserIDStr)

	// check sms verification code from user input
	code := r.FormValue("code")
	if code == "" {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrNoVerificationCode),
		})
		return
	}

	// check if user existed
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"_id": user.ID},
	}
	if err := mongoContext.QueryOne(sql, &user); err != nil {
		logrus.Warnf("query user by email error: %v", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrQueryUser),
		})
		return
	}
	if user.Verified {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserVerified),
		})
		return
	}

	// user not verified, send code
	sentCode := session.Values["verification_code"]
	if sentCode == nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrNoVerificationCode),
		})
		return
	}
	var sendCodeStr = sentCode.(string)
	if sendCodeStr != user.VerificationCode {
		logrus.Warnln("user verification code not equal the one got from session.")
		logrus.Warnln("overwrite verification code in session with user code.")
		session.Values["verification_code"] = user.VerificationCode
	}
	if !verification.Verify(&user, code) {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrVerificationCodeNotEqual),
		})
		return
	}

	updateSql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"_id": user.ID},
	}
	err = mongoContext.UpdateOne(updateSql, &user)
	if err != nil {
		logrus.Warnf("update user info error: %v", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserUpdate),
		})
		return
	}

	// verify user success, remove code in session
	delete(session.Values, "verification_code")
	session.Save(r, w)
	toJson(w, FormActionResponse{
		Error:   false,
		Message: locale(UserVerifiedSuccess),
	})
}
