package webservice

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/osinstorage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
)

type TokenTestSuit struct {
	suite.Suite
	service   *Service
	container *restful.Container
}

func (s *TokenTestSuit) SetupTest() {
	var err error

	client := &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://example.com",
	}

	osinStorage := osinstorage.NewMemoryStorage(*client)

	osinStorage.SaveAuthorize(&osin.AuthorizeData{
		RedirectUri: "http://example.com",
		Client:      client,
		CreatedAt:   time.Now(),
		ExpiresIn:   int32(10), // 10 seconds
		Code:        "auth_code",
	})

	osinServer := osin.NewServer(&config.DefaultConfig.OsinConfig, osinStorage)

	s.service, err = New(osinServer, nil)
	require.NoError(s.T(), err)

	s.container = restful.NewContainer()
	s.container.Add(s.service.WebService())
}

func TestTokenTestSuit(t *testing.T) {
	suite.Run(t, new(TokenTestSuit))
}

// As a valid client, I can request a token successfully with a authorize code.
func (s *TokenTestSuit) TestSignInWithValidUser() {

	// arrange: URL to send request
	u := url.URL{
		Scheme: "http",
		Path:   "/oauth2/token",
	}

	// arrange: from data for POST
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("client_id", "1234")
	data.Add("client_secret", "aabbccdd")
	data.Add("redirect_uri", "http://example.com")
	data.Add("code", "auth_code")

	// arrange: http request
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	require.NoError(s.T(), err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// action
	w := httptest.NewRecorder()
	s.container.ServeHTTP(w, req)

	// assert: http status code
	s.Equal(http.StatusOK, w.Code, w.Body.String())

	// assert: verify token dose exist
	gj := gjson.Parse(w.Body.String())
	s.Require().True(gj.Get("token_type").Exists(), w.Body.String())
	s.Equal("Bearer", gj.Get("token_type").String())
	s.Require().True(gj.Get("access_token").Exists(), w.Body.String())
	s.Regexp(".+", gj.Get("access_token").String())
	s.Require().True(gj.Get("refresh_token").Exists(), w.Body.String())
	s.Regexp(".+", gj.Get("refresh_token").String())
}
