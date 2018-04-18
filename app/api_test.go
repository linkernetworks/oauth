package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"bitbucket.org/linkernetworks/aurora/src/oauth/app/config"
	"bitbucket.org/linkernetworks/aurora/src/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/oauth/mongo"
	apptesting "bitbucket.org/linkernetworks/aurora/src/oauth/testing"
	"bitbucket.org/linkernetworks/aurora/src/oauth/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

const (
	oauthHost         = "0.0.0.0"
	oauthPort         = "9098"
	apiBaseUrl        = "http://0.0.0.0:9098"
	DefaultLocaleLang = "en-us"
	ZHLocaleLang      = "zh"
	TWLocaleLang      = "zh-tw"
)

var (
	mongoContext *mongo.Context
	mongoClient  *mongo.MongoClient
	appConfig    *config.AppConfig
)

func init() {
	configPath, err := apptesting.GetCurrentConfigPath()
	if err != nil {
		logrus.Fatal("config not found", err)
	}
	appConfig = config.Read(configPath)
	mongoClient = mongo.NewMongoClient(appConfig.MongoConfig)
	mongoContext = mongoClient.NewContext()

	// prepare test config
	gin.SetMode(gin.TestMode)
	clearDatabase()
}

func clearDatabase() {
	dbName := appConfig.MongoConfig.Database
	err := mongoContext.DropDatabase(dbName)
	if err != nil {
		logrus.Fatalf("delete mongo db %s error %v", dbName, err)
	} else {
		logrus.Infof("drop mongodb %s success", dbName)
	}
}

func PostJson(url string, data []byte) interface{} {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logrus.Fatalf("post app sinup error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("read response body error: %v", err)
	}

	var ret interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logrus.Fatalf("unmarshal response body error: %v", err)
	}

	return ret
}

// PostForm client post form encode data
func PostForm(url string, data url.Values) interface{} {
	resp, err := http.PostForm(url, data)
	if err != nil {
		logrus.Fatalf("post form data to %s error: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("read response body error: %v", err)
	}

	var ret interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logrus.Fatalf("unmarshal response body error: %v", err)
	}

	return ret
}

func PostFromWithCookie(url string, data url.Values, cookie *http.Cookie) interface{} {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		logrus.Fatalf("new post request error: %v\n", err)
	}
	req.AddCookie(cookie)
	// without below two lines, server could not parse code use gin(c.PosrForm("code"), code is "")
	// refer to https://groups.google.com/forum/#!topic/golang-nuts/x3JDvuWVO9A
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	req.Form = data

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatalf("get response from user verify error: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("read user verify response body error: %v\n", err)
	}

	var ret interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logrus.Fatalf("unmarshal response body error: %v\n", err)
	}

	return ret
}

func PostFormWithHeader(url string, data url.Values, headers map[string]string) interface{} {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		logrus.Fatalf("new post request error: %v\n", err)
	}

	// set header
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// without below two lines, server could not parse code use gin(c.PosrForm("code"), code is "")
	// refer to https://groups.google.com/forum/#!topic/golang-nuts/x3JDvuWVO9A
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	req.Form = data

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatalf("get response from user verify error: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("read user verify response body error: %v\n", err)
	}

	var ret interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logrus.Fatalf("unmarshal response body error: %v\n", err)
	}

	return ret
}

func setOAuthAuthorizeCookie(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["test@linker.com"] = "oauth.test.access.token"
	session.Save(r, w)
	toJson(w, nil)
}

func setUserVerifyCookie(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["current_user_id"] = verifyUserId
	session.Values["verification_code"] = verifyCode
	session.Save(r, w)
	toJson(w, nil)
}

func insertUser(user entity.User) {
	sql := mongo.Selector{
		Collection: "user",
	}
	err := mongoContext.InsertOne(sql, &user)
	if err != nil {
		logrus.Fatalf("insert test user to mongo error: %v", err)
	}
}

func upsertUser(user entity.User) {
	sql := mongo.Selector{
		Collection: "user",
	}
	err := mongoContext.Upsert(sql, &user)
	if err != nil {
		logrus.Fatalf("insert test user to mongo error: %v", err)
	}
}

func upsertApplication(app entity.Application) {
	sql := mongo.Selector{
		Collection: "application",
		Selector:   bson.M{"_id": app.ID},
	}
	err := mongoContext.Upsert(sql, &app)
	if err != nil {
		logrus.Fatalf("insert test app data to mongo error: %v", err)
	}
}

// testLocaleFunc return i18n.TranslateFunc for test
func testLocaleFunc(lang string) i18n.TranslateFunc {
	T, err := i18n.Tfunc(lang)
	if err != nil {
		logrus.Fatalf("must locale func error: %v\n", err)
	}

	return T
}

// helper funcs
func newTestApplication() entity.Application {
	app := entity.Application{
		ID:           bson.ObjectIdHex(clientId),
		Name:         "oauth-test",
		Description:  "Test001 for test",
		RedirectUri:  "http://localhost:14000/callback/code",
		ClientId:     clientId,
		ClientSecret: "oauth.test.client.secret",
		CreatedAt:    util.GetCurrentTimestamp(),
	}
	app.ClientId = clientId

	return app
}

func newTestUser() entity.User {
	u := entity.User{
		// ID:                    bson.NewObjectId(),
		ID:                    bson.ObjectId("111111111111"),
		Email:                 "test@linker.com",
		Password:              "123456",
		Cellphone:             "+8811111111111",
		Verified:              true,
		AccessToken:           "oauth.test.access.token",
		AccessTokenExpiryTime: util.GetCurrentTimestamp() + 3600,
	}

	return u
}

// refer to https://github.com/gin-gonic/gin/blob/master/gin_test.go#L24
func setupAPIServer(t *testing.T) {
	configPath, err := apptesting.GetCurrentConfigPath()
	if err != nil {
		logrus.Fatal("config not found", err)
	}

	go func(configPath string) {
		appService := NewServiceProviderFromConfig(appConfig)

		// warning: if use sessionStor after r.Route, server will panic:
		// Key "github.com/gin-contrib/sessions" does not exist
		store := sessions.NewCookieStore([]byte("something-very-secret"))
		if err != nil {
			logrus.Error(err)
		}
		store.Options = &sessions.Options{
			MaxAge: int(30 * time.Minute), //30min
			Path:   "/",
		}
		//testSession := store.Get(r, "oauth-test-session")

		// c, r := gin.CreateTestContext(httptest.NewRecorder())
		// r.Use(gin.Logger(), gin.Recovery())
		r := mux.NewRouter()
		/*
			r.Use(testSession)
			r.Use(cors.New(cors.Config{
				// AllowOrigins:     []string{"https://foo.com"},
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST"},
				AllowHeaders:     []string{"Origin"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			}))
		*/

		r.HandleFunc("/v1/me", CompositeServiceProvider(appService, handleMe))
		r.HandleFunc("/v1/email/check", CompositeServiceProvider(appService, checkEmailAvailability))
		r.HandleFunc("/v1/signup", CompositeServiceProvider(appService, UserSignUp))
		r.HandleFunc("/v1/signin", CompositeServiceProvider(appService, UserSignIn))
		r.HandleFunc("/v1/logout", CompositeServiceProvider(appService, UserLogout))
		r.HandleFunc("/v1/app/signup", CompositeServiceProvider(appService, AppSignup))
		r.HandleFunc("/oauth/pre", CompositeServiceProvider(appService, setOAuthAuthorizeCookie))
		r.HandleFunc("/v1/user/verify", CompositeServiceProvider(appService, UserVerify))
		r.HandleFunc("/v1/user/verify/request", CompositeServiceProvider(appService, setUserVerifyCookie))

		r.HandleFunc("/services/oauth/authorize", CompositeServiceProvider(appService, OAuthAuthorize))
		r.HandleFunc("/services/oauth/token", CompositeServiceProvider(appService, OAuthToken))

		// load i18n files, can be done in main func
		i18n.MustLoadTranslationFile("../config/locale/en-us.all.json")
		i18n.MustLoadTranslationFile("../config/locale/zh.all.json")
		i18n.MustLoadTranslationFile("../config/locale/zh-tw.all.json")

		http.ListenAndServe(":"+oauthPort, r)
		logrus.Infof("start api server success.")
	}(configPath)

	t.Log("waiting 2 second for server startup")
	time.Sleep(2 * time.Second)
}
