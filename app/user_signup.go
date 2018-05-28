package app

import (
	"fmt"
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/util"
	"bitbucket.org/linkernetworks/aurora/src/oauth/validator"
	"bitbucket.org/linkernetworks/aurora/src/oauth/verification"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// UserSignUp handle user signup
func UserSignUp(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	locale := MustLocaleFunc(r)

	var user entity.User
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Cellphone = r.FormValue("cellphone")
	// check user form data
	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(user.Email)
	if err != nil {
		validations["email"] = emailValidate
	}
	passworkValidate, err := validator.ValidatePassword(user.Password)
	if err != nil {
		validations["password"] = passworkValidate
	}
	phoneValidate, err := validator.ValidateCellphone(user.Cellphone)
	if err != nil {
		validations["cellphone"] = phoneValidate
	}
	if validations.HasError() {
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	// check if user existed
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"email": user.Email},
	}
	existed, err := mongoContext.CheckExist(sql, &user)
	if err != nil {
		logrus.Warnln("Check user existed error: ", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrCheckUserName),
		})
		return
	}
	if existed {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserExisted),
		})
		return
	}

	// user not existed in db, insert
	user.ID = bson.NewObjectId()

	salt := appService.OAuthConfig.Encryption.Salt
	encrypted, err := util.EncryptPassword(user.Password, salt)
	if err != nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}
	user.Password = encrypted
	if appService.SmsClient != nil {
		verification.Send(appService.SmsClient, &user, 6)
	}
	user.VerificationCode = verification.GenerateCode(4)
	user.CreatedAt = util.GetCurrentTimestamp()

	err = mongoContext.InsertOne(sql, &user)
	if err != nil {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrUserInsert),
		})
		return
	}

	// save verification code to session and send code to user
	msg := fmt.Sprintf(verification.SmsVerificationCodeMessageFormat, user.VerificationCode)
	_, except, err := appService.SmsClient.SendMessage(user.Cellphone, msg)
	if except != nil {
		logrus.Warnf("send sms verification code run into exception: %v\n", except)
	}
	if err != nil {
		logrus.Warnf("send sms verification code run into error: %v\n", err)
	}
	logrus.Infof("send verification sms code to user successfully.")

	// set user with verification code in session
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["verification_code"] = user.VerificationCode
	session.Values["current_user_id"] = user.ID.String()
	session.Save(r, w)

	http.Redirect(w, r, "/user/verify", http.StatusMovedPermanently)
}
