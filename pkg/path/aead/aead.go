package aead

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"

	"github.com/connesc/rlink/pkg/path"
)

type Authenticator struct {
	aead      cipher.AEAD
	nonceSize int
	overhead  int
}

func NewAuthenticator(aead cipher.AEAD) (*Authenticator, error) {
	return &Authenticator{
		aead:      aead,
		nonceSize: aead.NonceSize(),
		overhead:  aead.Overhead(),
	}, nil
}

func (r *Authenticator) FromOriginal(originalPath string) (string, error) {
	dir, file := path.Split(originalPath)

	payload := []byte(dir[:len(dir)-1])
	decodedPrefix := make([]byte, r.nonceSize, r.nonceSize+len(payload)+r.overhead)

	_, err := io.ReadFull(rand.Reader, decodedPrefix)
	if err != nil {
		return "", err
	}

	decodedPrefix = r.aead.Seal(decodedPrefix, decodedPrefix, payload, []byte(file))
	encodedPrefix := base64.RawURLEncoding.EncodeToString(decodedPrefix)

	return encodedPrefix + "/" + file, nil
}

var authenticatedPathPattern = regexp.MustCompile("^([^/]+)/([^/]*)$")

func (r *Authenticator) ToOriginal(authenticatedPath string) (string, error) {
	matches := authenticatedPathPattern.FindStringSubmatch(authenticatedPath)
	if matches == nil {
		return "", fmt.Errorf("aead: invalid authenticated path")
	}
	encodedPrefix := matches[1]
	file := matches[2]

	minPrefixSize := r.nonceSize + r.overhead
	prefixSize := base64.RawURLEncoding.DecodedLen(len(encodedPrefix))
	if prefixSize < minPrefixSize {
		return "", fmt.Errorf("aead: expected at least %v prefix bytes, got %v", minPrefixSize, prefixSize)
	}

	decodedPrefix := make([]byte, prefixSize)
	_, err := base64.RawURLEncoding.Decode(decodedPrefix, []byte(encodedPrefix))
	if err != nil {
		return "", fmt.Errorf("aead: invalid prefix encoding: %v", err)
	}

	payload := make([]byte, 0, prefixSize-r.nonceSize-r.overhead)
	payload, err = r.aead.Open(payload, decodedPrefix[:r.nonceSize], decodedPrefix[r.nonceSize:], []byte(file))
	if err != nil {
		return "", fmt.Errorf("aead: invalid prefix encryption: %v", err)
	}

	dir := string(payload) + "/"
	return path.Join(dir, file), nil
}
