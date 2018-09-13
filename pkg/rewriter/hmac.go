package rewriter

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"
)

type HMACRewriter struct {
	hashFunc func() hash.Hash
	key      []byte
}

func NewHMAC(hashFunc func() hash.Hash, key []byte) *HMACRewriter {
	return &HMACRewriter{
		hashFunc: hashFunc,
		key:      key,
	}
}

func (r *HMACRewriter) FromOriginal(originalPath string) (string, error) {
	trailingSlash := strings.HasSuffix(originalPath, "/")
	if trailingSlash {
		originalPath = originalPath[:len(originalPath)-1]
	}

	mac := hmac.New(r.hashFunc, r.key)
	mac.Write([]byte(originalPath))
	decodedMAC := mac.Sum(nil)
	encodedMAC := base64.RawURLEncoding.EncodeToString(decodedMAC)

	return join(encodedMAC, originalPath, trailingSlash), nil
}

func (r *HMACRewriter) ToOriginal(authPath string) (string, error) {
	trailingSlash := strings.HasSuffix(authPath, "/")
	if trailingSlash {
		authPath = authPath[:len(authPath)-1]
	}

	chunks := strings.SplitN(authPath, "/", 2)

	originalPath := ""
	switch len(chunks) {
	case 0:
		return "", fmt.Errorf("HMACRewriter: no MAC in authenticated path")
	case 2:
		originalPath = chunks[1]
	}
	encodedMAC := []byte(chunks[0])

	mac := hmac.New(r.hashFunc, r.key)
	expectedMACBytes := mac.Size()

	actualMACBytes := base64.RawURLEncoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("HMACRewriter: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err := base64.RawURLEncoding.Decode(decodedMAC, encodedMAC)
	if err != nil {
		return "", fmt.Errorf("HMACRewriter: invalid MAC: %v", err)
	}

	mac.Write([]byte(originalPath))
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("HMACRewriter: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	return join("", originalPath, trailingSlash), nil
}
