package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var tokenCallback = "http://localhost:14000/callback/token/code"

func startTokenCallback(t *testing.T) {
	// callback url likes
	// http://localhost:14000/callback/code?code=GdEbib0NTmiYR4yFa3pdrA&state=xyz
	go func() {
		http.HandleFunc("/callback/token/code", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			code := r.Form.Get("code")
			state := r.Form.Get("state")
			if code == "" || state == "" {
				t.Error()
			}

			u, err := url.Parse(apiBaseUrl + "/services/oauth/token")
			if err != nil {
				w.Write([]byte("internal error: " + err.Error()))
				return
			}

			p := u.Query()
			p.Set("grant_type", "authorization_code")
			p.Set("code", code)
			p.Set("redirect_uri", "")
			p.Set("client_id", clientId)
			p.Set("client_secret", "oauth.test.client.secret")
			p.Set("redirect_uri", tokenCallback)
			u.RawQuery = p.Encode()

			resp, err := http.Get(u.String())
			if err != nil {
				t.Errorf("get oauth token request error: %v", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("parse response body error: %v", err.Error)
			}

			var tokenRet map[string]interface{}
			err = json.Unmarshal(body, &tokenRet)
			if err != nil {
				t.Errorf("parse oauth token response error: %v", err)
			}

			/* sample when run into error
			map[
				error:server_error
				error_description:The authorization server encountered an unexpected condition that prevented it from fulfilling the request.
			]

			sample when success
			map[
				access_token:UTossiFySAGF1-XJAPoEsQ
				expires_in:3600
				refresh_token:VhP4yeyJR1mq7FqPf1hDVg
				token_type:Bearer
			]
			*/
			// error
			_, ok := tokenRet["error"]
			if ok {
				t.Errorf("access token error: %v", tokenRet["error_description"])
			}

			// success
			accessToken := tokenRet["access_token"].(string)
			if accessToken == "" {
				t.Errorf("access token should not be empty string")
			}
			refreshToken := tokenRet["refresh_token"].(string)
			if refreshToken == "" {
				t.Errorf("refresh token should not be empty string")
			}
		})

		http.ListenAndServe(":14000", nil)
	}()

	t.Log("waiting 2 second for oauth authorize callback server startup")
	time.Sleep(2 * time.Second)
}

func TestOauthToken(t *testing.T) {
	app := newTestApplication()
	app.RedirectUri = tokenCallback
	upsertApplication(app) // from oauth_authorize_test.go
	user := newTestUser()
	insertUser(user)
	startTokenCallback(t)
	setupAPIServer(t)

	authorizeCodeUrl, err := url.Parse(apiBaseUrl + "/services/oauth/authorize")
	if err != nil {
		t.Errorf("parse authorize url error: %v", err)
	}

	requestValues := url.Values{}
	requestValues.Set("response_type", "code")
	requestValues.Set("client_id", "4d88e15b60f486e428412dc9")
	requestValues.Set("state", "xyz")
	requestValues.Set("redirect_uri", url.QueryEscape(tokenCallback))
	requestValues.Set("user_access_token", "oauth.test.access.token")
	requestValues.Set("user_email", user.Email)

	authorizeCodeUrl.RawQuery = requestValues.Encode()

	resp, err := http.Get(authorizeCodeUrl.String())
	if err != nil {
		t.Errorf("get request to authorize error: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("TestUserSignIn done, clear mongo db")
	clearDatabase()
}
