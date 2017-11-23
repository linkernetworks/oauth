package app

import (
	"net/http"

	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/entity"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/mongo"
	"bitbucket.org/linkernetworks/cv-tracker/src/oauth/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func checkEmailAvailability(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()

	var user entity.User
	user.Email = r.FormValue("email")

	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(user.Email)
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

	// check if user existed
	sql := mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"email": user.Email},
	}
	existed, err := mongoContext.CheckExist(sql, &user)
	if err != nil {
		logrus.Warnf("mongo error: %v", err)
		toJson(w, gin.H{
			"error":   true,
			"message": "Internal Server Error",
		})
		return
	}

	if !existed {
		validations["email"] = validator.FieldValidation{Field: "email", Error: true, Message: "This email doesn't exist."}
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	toJson(w, gin.H{
		"error":   false,
		"message": "email OK",
	})
}

func loadUserByEmail(c *mongo.MongoClient, email string) (*entity.User, error) {
	mongoContext := c.NewContext()
	defer mongoContext.Close()

	var sql = mongo.Selector{
		Collection: "user",
		Selector:   bson.M{"email": email},
	}

	var u entity.User
	existed, err := mongoContext.CheckExist(sql, &u)
	if existed {
		return &u, err
	}
	return nil, err
}
