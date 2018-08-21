package rewriter

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"
)

type PathSigner struct {
	hashFunc func() hash.Hash
	key      []byte
	encoding *base64.Encoding
}

func NewPathSigner(hashFunc func() hash.Hash, key []byte, encoding *base64.Encoding) *PathSigner {
	return &PathSigner{
		hashFunc: hashFunc,
		key:      key,
		encoding: encoding,
	}
}

func (s *PathSigner) FromOriginal(originalPath string) (string, error) {
	mac := hmac.New(s.hashFunc, s.key)
	mac.Write([]byte(originalPath))
	decodedMAC := mac.Sum(nil)
	encodedMAC := s.encoding.EncodeToString(decodedMAC)

	return encodedMAC + originalPath, nil
}

func (s *PathSigner) ToOriginal(signedPath string) (string, error) {
	chunks := strings.SplitN(signedPath[1:], "/", 2)

	originalPath := ""
	switch len(chunks) {
	case 0:
		return "", fmt.Errorf("PathSigner: no MAC in signed path")
	case 2:
		originalPath = chunks[1]
	}
	encodedMAC := []byte(chunks[0])

	mac := hmac.New(s.hashFunc, s.key)
	expectedMACBytes := mac.Size()

	actualMACBytes := s.encoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("PathSigner: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err := s.encoding.Decode(decodedMAC, encodedMAC)
	if err != nil {
		return "", fmt.Errorf("PathSigner: invalid MAC: %v", err)
	}

	mac.Write([]byte(originalPath))
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("PathSigner: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	return originalPath, nil

}
