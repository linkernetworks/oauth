package storage

import (
	"time"

	"github.com/linkernetworks/oauth/src/entity"
)

type AuthorizeData struct {
	// Client information
	// replace osin.Client to avoid mgo parse error
	// refer to https://github.com/RangelReale/osin/issues/40
	Client entity.Application

	// Authorization code
	Code string

	// Token expiration in seconds
	ExpiresIn int32

	// Requested scope
	Scope string

	// Redirect Uri from request
	RedirectUri string

	// State data from request
	State string

	// Date created
	CreatedAt time.Time

	// Data to be passed to mongo. Not used by the library.
	UserData interface{}

	// Optional code_challenge as described in rfc7636
	CodeChallenge string
	// Optional code_challenge_method as described in rfc7636
	CodeChallengeMethod string
}
