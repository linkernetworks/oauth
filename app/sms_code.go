package app

import (
	"fmt"
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/verification"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// SmsCode send verification sms code when registering device
func SmsCode(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	var user entity.User

	// because access memory is more quick than mongodb
	// check user id in session first, then mongodb
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := session.Values["current_user_id"]
	if userID == nil {
		toJson(w, gin.H{
			"message": "send sms code error, not provide user id in cookie.",
		})
		return
	}
	var newUserID string = userID.(string)

	userSMSCode := session.Values["verification_code"]
	// userSMSCode is not nil, we have sent sms code
	if userSMSCode != nil {
		newUserSMSCode := userSMSCode.(string)
		toJson(w, gin.H{
			"message":  "We have sent sms code to you",
			"sms_code": newUserSMSCode,
		})
		return
	}

	// userSMSCode is nil, check userId and send code
	user.ID = bson.ObjectId(newUserID)
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"_id": user.ID},
	}
	exited, err := mongoContext.CheckExist(sql, &user)
	if err != nil {
		logrus.Warnf("check user exited error: %v", err)
		toJson(w, gin.H{
			"message": "check user exited error.",
		})
		return
	}
	if !exited {
		logrus.Warnf("user with id %v does not existed in db.\n", user.ID)
		toJson(w, gin.H{
			"message": fmt.Sprintf("user with id %v does not existed in db", user.ID),
		})
		return
	}

	// user with id existed, generate verification code and update db
	user.VerificationCode = verification.GenerateCode(4)
	sql = mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"_id": user.ID},
	}
	err = mongoContext.UpdateOne(sql, &user)
	if err != nil {
		logrus.Warnln("update user with verification code error, please retry.")
		toJson(w, gin.H{
			"message": "update user with verification code error, please retry.",
		})
		return
	}

	// send code to user device
	msg := fmt.Sprintf(verification.SmsVerificationCodeMessageFormat, user.VerificationCode)
	_, except, err := appService.SmsClient.SendMessage(user.Cellphone, msg)
	if except != nil {
		logrus.Warnf("send sms code run into exception: %v", except)
		toJson(w, gin.H{
			"message": "send sms code run into exception.",
		})
		return
	}
	if err != nil {
		logrus.Warnf("send sms code run into error: %v", err)
		toJson(w, gin.H{
			"message": "send sms code run into error.",
		})
		return
	}

	// set user sms code in session
	session.Values["verification_code"] = user.VerificationCode
	session.Save(r, w)

	toJson(w, gin.H{
		"message":  "send sms code successfully.",
		"sms_code": user.VerificationCode,
	})
}
