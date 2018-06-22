package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

// EncryptPassword encrypt password
func EncryptPassword(password, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return string(dk), nil
}

// GetCurrentTimestamp get timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GenPubPriKey gen private and public key by rsa algorithm
func GenPubPriKey(bitSize int) (pri, pub string) {
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		logrus.Fatalf("generate key error: %s", err)
	}

	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	privateKey := pem.EncodeToMemory(privateBlock)
	privateKeyStr := string(privateKey)
	pri = removeKeyExtra(privateKeyStr)

	asn1Bytes, err := asn1.Marshal(key.PublicKey)
	if err != nil {
		logrus.Fatalf("marshal publickey error:%s", err)
	}
	var pemBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	pemKey := pem.EncodeToMemory(pemBlock)
	pemKetStr := string(pemKey)
	pub = removeKeyExtra(pemKetStr)

	return
}

// removeKeyExtra key has three lines, remove key first and third line
func removeKeyExtra(key string) string {
	keyArray := strings.Split(key, "\n")
	return keyArray[1]
}
