package rewriter

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"
)

type URLSigner struct {
	hashFunc func() hash.Hash
	key      []byte
	encoding *base64.Encoding
}

func NewURLSigner(hashFunc func() hash.Hash, key []byte, encoding *base64.Encoding) *URLSigner {
	return &URLSigner{
		hashFunc: hashFunc,
		key:      key,
		encoding: encoding,
	}
}

func (s *URLSigner) FromOriginal(originalPath string) (string, error) {
	mac := hmac.New(s.hashFunc, s.key)
	mac.Write([]byte(originalPath))
	decodedMAC := mac.Sum(nil)
	encodedMAC := s.encoding.EncodeToString(decodedMAC)

	return encodedMAC + originalPath, nil
}

func (s *URLSigner) ToOriginal(signedPath string) (string, error) {
	chunks := strings.SplitN(signedPath[1:], "/", 2)

	originalPath := ""
	switch len(chunks) {
	case 0:
		return "", fmt.Errorf("URLSigner: no MAC in request URL")
	case 2:
		originalPath = chunks[1]
	}
	encodedMAC := []byte(chunks[0])

	mac := hmac.New(s.hashFunc, s.key)
	expectedMACBytes := mac.Size()

	actualMACBytes := s.encoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("URLSigner: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err := s.encoding.Decode(decodedMAC, encodedMAC)
	if err != nil {
		return "", fmt.Errorf("URLSigner: invalid MAC: %v", err)
	}

	mac.Write([]byte(originalPath))
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("URLSigner: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	return originalPath, nil

}
