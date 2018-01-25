package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/linkernetworks/aurora/src/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	"bitbucket.org/linkernetworks/aurora/src/oauth/util"
	"github.com/RangelReale/osin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoContext *mongo.Context
	mongoClient  *mongo.MongoClient
	osinStorage  *OsinStorage
	configPath   string
	clientId     string
)

func init() {
	cp, found := os.LookupEnv("OAUTH_CONFIG_PATH")
	if !found {
		p, err := filepath.Abs("../../mongo/test.json")
		if err != nil {
			log.Fatal(err)
		}
		configPath = p
	} else {
		configPath = cp
	}
	mongoConfig := newMongoConfig()
	mongoClient = mongo.NewMongoClient(mongoConfig)
	clientId = "4d88e15b60f486e428412dc9"
	mongoContext = mongoClient.NewContext()
	osinStorage = NewOsinStorage(mongoConfig)
}

func newMongoConfig() mongo.MongoConfig {
	configPath = os.Getenv("OAUTH_CONFIG_PATH")
	if configPath == "" {
		configPath = "../../mongo/test.json"
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("read config file %s error: %v\n", configPath, err)
	}

	var mongoConfig mongo.MongoConfig
	err = json.Unmarshal(content, &mongoConfig)
	if err != nil {
		logrus.Fatalf("parse mongo config error: %v\n", err)
	}

	return mongoConfig
}

// helper funcs
func newTestApplication() entity.Application {
	app := entity.Application{
		ID:          bson.ObjectIdHex(clientId),
		Name:        "Test001",
		Description: "Test001 for test",
		RedirectUri: "http://localhost:9090/callback",
		CreatedAt:   util.GetCurrentTimestamp(),
	}
	app.GenerateToken()
	app.ClientId = clientId

	return app
}

func newTestAuthorizeData() osin.AuthorizeData {
	authorizeData := osin.AuthorizeData{
		Client:      newTestApplication(),
		Code:        "1234",
		Scope:       "all",
		RedirectUri: "http://localhost:9090/callback",
		State:       "xyz",
	}

	return authorizeData
}

func newTestAccessData() osin.AccessData {
	authorizeData := newTestAuthorizeData()
	accessData := osin.AccessData{
		Client:        newTestApplication(),
		AuthorizeData: &authorizeData,
		AccessToken:   "access_token",
		RefreshToken:  "refresh_token",
		Scope:         "all",
		RedirectUri:   "http:localhost:9090/callback",
	}

	return accessData
}

func insertSampleClient() {
	app := newTestApplication()
	sql := mongo.Selector{
		Collection: CLIENT_COLLECTION,
	}

	// check if app existed
	existed, err := mongoContext.CheckExist(sql, &app)
	if err != nil {
		log.Fatalf("check app existed error: %v", err)
	}
	if existed {
		// app name has existed. skip
		log.Printf("app has existed, skip\n")
		return
	}

	err = mongoContext.InsertOne(sql, &app)
	if err != nil {
		log.Fatalf("insert sample app to mongo error: %v", err)
	}
}

func deleteSampleClient() {
	sql := mongo.Selector{
		Collection: CLIENT_COLLECTION,
		Selector:   bson.M{"name": "Test001"},
	}

	err := mongoContext.RemoveOne(sql)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			log.Printf("doc to delete does not found, skip")
			return
		}
		log.Fatalf("insert sample app to mongo error: %v", err)
	}
}

func TestClone(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	newm := osinStorage.Clone()
	assert.NotNil(t, newm)
}

func TestGetClient(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	insertSampleClient()
	defer deleteSampleClient()

	client, err := osinStorage.GetClient(clientId)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("not found client with name %s from mongo, please insert first \n", clientId)
			return
		}
		t.Errorf("get client by id %s error: %v", clientId, err)
	}
	assert.Equal(t, client.GetId(), clientId)
}

func TestSaveAuthorize(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	authorizeData := newTestAuthorizeData()
	err := osinStorage.SaveAuthorize(&authorizeData)
	if err != nil {
		t.Errorf("SaveAuthorize error: %v", err)
	}
}

func TestLoadAuthorize(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	_, err := osinStorage.LoadAuthorize("1234")
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("not found aothurizeData, please insert first")
			return
		}
		t.Errorf("locad authorize data error: %v", err)
	}
}

func TestRemoveAuthorize(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	TestSaveAuthorize(t)
	code := "1234"
	err := osinStorage.RemoveAuthorize(code)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("doc with code %s does not found, please insert first", code)
			return
		}

		t.Errorf("remove authorize data error: %v", err)
	}
}

func TestSaveAccess(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	accessData := newTestAccessData()
	err := osinStorage.SaveAccess(&accessData)
	if err != nil {
		t.Errorf("save access data to mongo error: %v", err)
	}
}

func TestLoadAccess(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	token := "access_token"
	_, err := osinStorage.LoadAccess(token)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("loadAccess with token %s does not found, please insert first", token)
			return
		}

		t.Errorf("load access data error: %v", err)
	}
}

func TestRemoveAccess(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	token := "access_token"
	err := osinStorage.RemoveAccess(token)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("remove access data with token %s does not found, please insert first", token)
			return
		}

		t.Errorf("remove access data error: %v", err)
	}
}

func TestLoadRefresh(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	token := "access_token"
	_, err := osinStorage.LoadRefresh(token)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("loadRefresh with token %s does not found, please insert first", token)
			return
		}

		t.Errorf("load fresh data error: %v", err)
	}
}

func TestRemoveRefresh(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_K8S"); !ok {
		t.Skip("Skip kubernetes related tests")
		return
	}

	token := "access_token"
	err := osinStorage.RemoveRefresh(token)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			t.Skipf("removeRefresh with token %s does not found, please insert first", token)
			return
		}

		t.Errorf("remove fresh error: %v", err)
	}
}
