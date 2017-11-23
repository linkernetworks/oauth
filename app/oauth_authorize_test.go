package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const clientId = "4d88e15b60f486e428412dc9"

var authorizeCallback = "http://localhost:14000/callback/authorize/code"

func startClientCallback(t *testing.T) {
	go func() {
		http.HandleFunc("/callback/authorize/code", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			code := r.Form.Get("code")
			state := r.Form.Get("state")
			ret := make(map[string]string)
			ret["state"] = state
			ret["code"] = code
			retJson, err := json.Marshal(&ret)
			if err != nil {
				logrus.Fatalf("marshal map to json error: %v", err)
			}
			w.Header().Set("Content-type", "application/json")
			w.Write(retJson)
		})

		http.ListenAndServe(":14000", nil)
	}()

	t.Log("waiting 2 second for oauth authorize callback server startup")
	time.Sleep(2 * time.Second)
}

func getOAuthAuthorizeCookie() *http.Cookie {
	resp, err := http.Get(apiBaseUrl + "/oauth/pre")
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Fatalf("pre session request error")
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		logrus.Fatalf("cookies length expected to be 1")
	}

	return cookies[0]
}

func TestOAuthAuthorizeCode(t *testing.T) {
	app := newTestApplication()
	app.RedirectUri = authorizeCallback
	upsertApplication(app)
	user := newTestUser()
	insertUser(user)
	startClientCallback(t)
	setupAPIServer(t)
	cookie := getOAuthAuthorizeCookie()

	authorizeCodeUrl, err := url.Parse(apiBaseUrl + "/services/oauth/authorize")
	if err != nil {
		t.Errorf("parse authorize url error: %v", err)
	}

	requestValues := url.Values{}
	requestValues.Set("response_type", "code")
	requestValues.Set("client_id", clientId)
	requestValues.Set("state", "xyz")
	requestValues.Set("redirect_uri", url.QueryEscape(authorizeCallback))
	requestValues.Set("user_email", "test@linker.com")

	authorizeCodeUrl.RawQuery = requestValues.Encode()

	req, err := http.NewRequest("GET", authorizeCodeUrl.String(), nil)
	if err != nil {
		t.Errorf("build request of url %s error: %v", authorizeCodeUrl.String(), err)
	}
	req.AddCookie(cookie)
	client := &http.Client{}
	// resp, err := http.Get(authorizeCodeUrl.String())
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("get request to authorize error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("parse authorize response body error: %v", err)
	}
	var authorizeResp map[string]string
	err = json.Unmarshal(body, &authorizeResp)
	if err != nil {
		t.Errorf("parse authorize response body to json error: %v", err)
	}

	assert.Equal(t, authorizeResp["state"], "xyz")
	if len(authorizeResp["code"]) <= 0 {
		t.Error()
	}

	t.Logf("TestUserSignIn done, clear mongo db")
	clearDatabase()
}
