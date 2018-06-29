package osinstorage

import (
	"testing"

	"github.com/RangelReale/osin"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {

	// arrange
	expected := &MemoryStorage{
		clients: map[string]osin.DefaultClient{
			"client_id": osin.DefaultClient{
				Id: "client_id",
			},
		},
		authorize: map[string]osin.AuthorizeData{
			"client_id": osin.AuthorizeData{
				Code: "codeeeee",
			},
		},
		access: map[string]osin.AccessData{
			"client_id": osin.AccessData{
				AccessToken: "token1",
			},
		},
		refresh: map[string]string{
			"token2": "token1",
		},
	}

	// action
	actual := expected.Clone()

	// assert
	assert.Equal(t, expected, actual)
}

func TestGetAndSaveAccess(t *testing.T) {
	// arrange
	storage := NewMemoryStorage()
	expectedAccess := &osin.AccessData{
		AuthorizeData: &osin.AuthorizeData{
			Code: "codeeeee",
		},
		AccessToken: "toeknnnnn",
	}

	// action
	storage.SaveAccess(expectedAccess)
	actualAccess, err := storage.LoadAccess("codeeeee")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAccess, actualAccess)
}
