package hmac

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"regexp"
)

type Authenticator struct {
	hashFunc func() hash.Hash
	key      []byte
}

func NewAuthenticator(hashFunc func() hash.Hash, key []byte) *Authenticator {
	return &Authenticator{
		hashFunc: hashFunc,
		key:      key,
	}
}

func (r *Authenticator) FromOriginal(originalPath string) (string, error) {
	payload := originalPath
	if payload == "/" {
		payload = ""
	}

	mac := hmac.New(r.hashFunc, r.key)
	mac.Write([]byte(payload))
	decodedMAC := mac.Sum(nil)
	encodedMAC := base64.RawURLEncoding.EncodeToString(decodedMAC)

	return encodedMAC + "/" + payload, nil
}

var authenticatedPathPattern = regexp.MustCompile("^([^/]+)/(.*)$")

func (r *Authenticator) ToOriginal(authenticatedPath string) (string, error) {
	matches := authenticatedPathPattern.FindStringSubmatch(authenticatedPath)
	if matches == nil {
		return "", fmt.Errorf("hmac: invalid authenticated path")
	}
	encodedMAC := matches[1]
	payload := matches[2]

	mac := hmac.New(r.hashFunc, r.key)
	expectedMACBytes := mac.Size()

	actualMACBytes := base64.RawURLEncoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("hmac: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err := base64.RawURLEncoding.Decode(decodedMAC, []byte(encodedMAC))
	if err != nil {
		return "", fmt.Errorf("hmac: invalid MAC: %v", err)
	}

	mac.Write([]byte(payload))
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("hmac: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	originalPath := payload
	if originalPath == "" {
		originalPath = "/"
	}
	return originalPath, nil
}
