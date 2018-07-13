package httphandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/linkernetworks/oauth/src/osinstorage"

	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/stretchr/testify/assert"
)

// As an OAuth2 client with valid authorized code,
// 'access_token' & 'refresh_token' tokens should return from the endpoint /token.
func TestToken(t *testing.T) {
	// arrange: prepare static data
	clientID := "client_id"
	clientSecret := "client_secret"
	redirectURI := "http://aaa.bbb"
	oauthCode := "oauth_code"
	client := &osin.DefaultClient{
		Id:          clientID,
		Secret:      clientSecret,
		RedirectUri: redirectURI,
	}

	// arrange: Storage for osin
	osinStorage := osinstorage.NewMemoryStorage(*client)
	osinStorage.SaveAuthorize(&osin.AuthorizeData{
		RedirectUri: redirectURI,
		Client:      client,
		CreatedAt:   time.Now(),
		ExpiresIn:   int32(10), // 10 seconds
		Code:        oauthCode,
	})

	// arrange: Create osin server
	osinServer := osin.NewServer(&config.DefaultConfig.OsinConfig, osinStorage)

	// arrange: Create dummy router
	router := gin.New()
	router.POST("/token", func(c *gin.Context) {
		c.Set("osinServer", osinServer)
		Token(c)
	})

	// arrange: prepare HTTP request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("redirect_uri", redirectURI)
	data.Add("code", oauthCode)
	req, _ := http.NewRequest("POST", "/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// action
	router.ServeHTTP(w, req)
	var actual map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &actual)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, actual, "access_token")
	assert.Contains(t, actual, "refresh_token")
}
