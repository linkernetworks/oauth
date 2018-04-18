package app

import (
	_ "github.com/gin-contrib/cors"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/sirupsen/logrus"
	"net/http"
	"bitbucket.org/linkernetworks/aurora/src/service/provider"
)

const (
	CERT_PRIVATE_KEY = "./tls-key/server.key"
	CERT_PUBLIC_KEY  = "./tls-key/server.crt"
)

type RequestHandler func(w http.ResponseWriter, r *http.Request, as *provider.Container)

func AppRoutes(appService *provider.Container) *mux.Router {
	// warning: if use sessionStor after r.Route, server will panic:
	// Key "github.com/gin-contrib/sessions" does not exist
	r := mux.NewRouter()
	/*
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
	r.HandleFunc("/v1/code/resend", CompositeServiceProvider(appService, SmsCode))
	r.HandleFunc("/v1/user/verify", CompositeServiceProvider(appService, UserVerify))

	r.HandleFunc("/services/oauth/authorize", CompositeServiceProvider(appService, OAuthAuthorize))
	r.HandleFunc("/services/oauth/token", CompositeServiceProvider(appService, OAuthToken))

	//	r.LoadHTMLGlob("./public/templates/*")
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/static")))
	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./public/assets")))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/robots.txt")
	})
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/signin.html")
	})
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/signup.html")
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/logout.html")
	})
	r.HandleFunc("/user/verify", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/user_verify.html")
	})

	// load i18n files, can be done in main func
	i18n.MustLoadTranslationFile("./config/locale/en-us.all.json")
	i18n.MustLoadTranslationFile("./config/locale/zh.all.json")
	i18n.MustLoadTranslationFile("./config/locale/zh-tw.all.json")

	return r
}

func Start(bind string, appService *provider.Container) error {
	return http.ListenAndServe(bind, AppRoutes(appService))
}

func StartSsl(bind string, appService *provider.Container) error {
	return http.ListenAndServeTLS(bind, CERT_PUBLIC_KEY, CERT_PRIVATE_KEY, AppRoutes(appService))
}

// CompositeServiceProvider apply mongo client to HandlerFunc
func CompositeServiceProvider(appService *provider.Container, handler RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, appService)
	}
}

// MustLocaleFunc return i18n.TranslateFunc in each request
func MustLocaleFunc(r *http.Request) i18n.TranslateFunc {
	acceptLang := r.Header.Get("Accept-Language")
	defaultLang := "en-us"
	T, err := i18n.Tfunc(acceptLang, defaultLang)
	if err != nil {
		logrus.Fatalf("must locale func error: %v\n", err)
	}

	return T
}
