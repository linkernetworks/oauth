package webservice

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/RangelReale/osin"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/osinstorage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthorizeTestSuit struct {
	suite.Suite
	service   *Service
	container *restful.Container
}

func (s *AuthorizeTestSuit) SetupTest() {
	var err error

	osinStorage := osinstorage.NewMemoryStorage(
		osin.DefaultClient{
			Id:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://example.com",
		},
	)

	osinServer := osin.NewServer(&config.DefaultConfig.OsinConfig, osinStorage)

	s.service, err = New(osinServer, nil)
	require.NoError(s.T(), err)

	s.container = restful.NewContainer()
	s.container.Add(s.service.WebService())
}

func TestAuthorizeTestSuit(t *testing.T) {
	suite.Run(t, new(AuthorizeTestSuit))
}

// With a valid client_id, a redirect responce to client redirect_url should be returned from /oauth2/authorize.
func (s *AuthorizeTestSuit) TestSignInWithValidUser() {
	// arrange
	u := url.URL{
		Scheme:   "http",
		Path:     "/authorize",
		RawQuery: "client_id=1234&response_type=code",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	require.NoError(s.T(), err)

	// action
	w := httptest.NewRecorder()
	s.container.ServeHTTP(w, req)

	// assert
	s.Equal(http.StatusFound, w.Code)
	s.Regexp("^(http://example.com\\?).*(code=).+$", w.Header().Get("Location"))
}
