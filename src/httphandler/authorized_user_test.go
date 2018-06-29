package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

// Test func AuthorizedUserOrRedirect()
// As a valid user, user should get status code 200.
func TestAuthorizedUserOrRedirect(t *testing.T) {
	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/test_path", nil)
	w := httptest.NewRecorder()

	// arrange: store auth data in session store
	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)
	session, _ := store.Get(req, "session")
	session.Values["user_id"] = "aa"
	session.Save(req, w)

	// arrange: router
	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.POST("/test_path", AuthorizedUserOrRedirect)

	// action
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
}

// Test func AuthorizedUserOrRedirect()
// As a invalid user, user should be redirected to login page.
func TestAuthorizedUserOrRedirectRedirectUser(t *testing.T) {
	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/test_path", nil)
	req.RequestURI = "redirect_uriiiii"
	w := httptest.NewRecorder()

	// arrange: store auth data in session store
	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)

	// arrange: router
	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.POST("/test_path", AuthorizedUserOrRedirect)

	// action
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 307, w.Code)
	assert.Equal(t, "/login?redirect_uri=redirect_uriiiii", w.Header().Get("Location"))
}

// Test func AuthorizedUser()
// As a valid user, user should get status code 200.
func TestAuthorizedUser(t *testing.T) {
	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/test_path", nil)
	w := httptest.NewRecorder()

	// arrange: store auth data in session store
	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)
	session, _ := store.Get(req, "session")
	session.Values["user_id"] = "aa"
	session.Save(req, w)

	// arrange: router
	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.POST("/test_path", AuthorizedUser)

	// action
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
}

// Test func AuthorizedUser()
// As a invalid user, user should get status code 401.
func TestAuthorizedUserWithInvalidUser(t *testing.T) {
	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/test_path", nil)
	w := httptest.NewRecorder()

	// arrange: store auth data in session store
	secret := securecookie.GenerateRandomKey(64)
	store := memstore.NewStore(secret)

	// arrange: router
	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.POST("/test_path", AuthorizedUser)

	// action
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 401, w.Code)
}
