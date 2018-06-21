package app

import (
	"net/http"

	"github.com/linkernetworks/oauth/src/entity"
	"github.com/linkernetworks/oauth/src/mongo"
	"github.com/linkernetworks/oauth/src/util"
	"github.com/linkernetworks/oauth/src/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// AppSignup handler app signup
func AppSignup(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	mongoContext := appService.MongoClient.NewContext()
	defer mongoContext.Close()
	var app entity.Application
	var validations = validator.ValidationMap{}
	locale := MustLocaleFunc(r)

	app.Name = r.FormValue("name")
	if app.Name == "" {
		validations["name"] = validator.FieldValidation{Field: "name", Error: true, Message: "name is required."}
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	// check if app name exist
	sql := mongo.Selector{
		Collection: "application",
		Selector:   bson.M{"name": app.Name},
	}

	existed, err := mongoContext.CheckExist(sql, &app)
	if err != nil {
		logrus.Warnf("check app name existed error: %v", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrCheckAppName),
		})
		return
	}
	if existed {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrAppNameExisted),
		})
		return
	}

	app.ID = bson.NewObjectId()
	app.GenerateToken()
	app.CreatedAt = util.GetCurrentTimestamp()

	// save to db
	sql = mongo.Selector{
		Collection: "application",
	}
	err = mongoContext.InsertOne(sql, &app)
	if err != nil {
		logrus.Warnf("save app into mongo error: %v", err)
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrSaveApp),
		})
		return
	}

	toJson(w, gin.H{
		"message":       locale(SaveAppSuccess),
		"name":          app.Name,
		"client_id":     app.ClientId,
		"client_secret": app.ClientSecret,
	})
}
