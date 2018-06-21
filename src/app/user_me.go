package app

import (
	"net/http"
	"time"

	"github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/oauth/mongo"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func handleMe(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()

	token := r.FormValue("access_token")

	if token == "" {
		logrus.Warnf("Empty token")
		toJson(w, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	var user entity.User
	user.AccessToken = token

	// check if user existed
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"token": user.AccessToken},
	}
	existed, err := mongoContext.CheckExist(sql, &user)
	if err != nil || !existed {
		logrus.Warnf("check user existed error: %v", err)
		toJson(w, gin.H{
			"message": "Invalid token",
			"error":   true,
		})
		return
	}

	// check if token is expired
	currentTime := time.Now().Unix()
	if user.AccessTokenExpiryTime < currentTime {
		logrus.Warnf("user token is expired: ", user.Email)
		toJson(w, gin.H{
			"message": "Expired token",
			"error":   true,
		})
		return
	}

	toJson(w, gin.H{
		"message":  "me success.",
		"UserInfo": user,
	})
}
