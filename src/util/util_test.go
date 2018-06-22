package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptPassword(t *testing.T) {
	s, err := EncryptPassword("test", "12345678")
	assert.NoError(t, err)
	if len(s) <= 0 {
		t.Errorf("test %s error", "EncryptPas")
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	ts := GetCurrentTimestamp()
	if ts <= 0 {
		t.Errorf("gen timestamp error")
	}
}

func TestGenPubPriKey(t *testing.T) {
	size := 16
	pri, pub := GenPubPriKey(size)
	if len(pri) <= 0 || len(pub) != size {
		t.Errorf("gen public and private key error")
	}
}

func TestRemoveKeyExtra(t *testing.T) {
	s := `-----BEGIN PRIVATE KEY-----
MCgCAQACBAC2jV0CAwEAAQIEAItOAQICDrECAgxtAgIAwQICCzkCAgQU
-----END PRIVATE KEY-----`
	expected := "MCgCAQACBAC2jV0CAwEAAQIEAItOAQICDrECAgxtAgIAwQICCzkCAgQU"
	ret := removeKeyExtra(s)
	if ret != expected {
		t.Errorf("remove extra error, expected %s but got %s", expected, ret)
	}
}
