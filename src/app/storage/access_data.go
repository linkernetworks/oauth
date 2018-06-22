package storage

import (
	"time"

	"github.com/linkernetworks/oauth/src/entity"
)

// AccessData represents an access grant (tokens, expiration, client, etc)
type AccessData struct {
	// Client information
	Client entity.Application

	// Authorize data, for authorization code
	AuthorizeData *AuthorizeData

	// Previous access data, for refresh token
	AccessData *AccessData

	// Access token
	AccessToken string

	// Refresh Token. Can be blank
	RefreshToken string

	// Token expiration in seconds
	ExpiresIn int32

	// Requested scope
	Scope string

	// Redirect Uri from request
	RedirectUri string

	// Date created
	CreatedAt time.Time

	// Data to be passed to mongo. Not used by the library.
	UserData interface{}
}
