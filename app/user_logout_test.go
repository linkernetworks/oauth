package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogout(t *testing.T) {

	// EXECUTOR_NUMBER is a jenkins environment variable
	if os.Getenv("EXECUTOR_NUMBER") != "" {
		t.Skip("Fix this for concurrent build")
	}

	// set cookie for test
	// use func in aouth_authorize_test.go
	setupAPIServer(t)
	getOAuthAuthorizeCookie()
	logoutUrl := apiBaseUrl + "/v1/logout"

	// test success, redirect to /signin, /signin not exists in route
	// so http code is 404
	resp1, err := http.Get(logoutUrl + "?user_email=test@linker.com")
	if err != nil {
		t.Errorf("request logout error: %v", err)
	}
	defer resp1.Body.Close()
	assert.Equal(t, resp1.StatusCode, http.StatusNotFound)

	// test logout error
	// no email
	resp2, err := http.Get(logoutUrl)
	if err != nil {
		t.Errorf("request logout error: %v", err)
	}
	defer resp1.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Errorf("read logout response error: %v", err)
	}

	var ret2 map[string]interface{}
	err = json.Unmarshal(body2, &ret2)
	if err != nil {
		t.Errorf("parse logout response body to map error: %v", err)
	}

	assert.Equal(t, ret2["error"].(bool), true)
}
