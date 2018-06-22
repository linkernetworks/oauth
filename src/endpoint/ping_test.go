package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	// arrange
	e := New(nil)

	// assert
	assert.HTTPBodyContains(t, e.Ping, "GET", "/ping", nil, "Pong")
}
