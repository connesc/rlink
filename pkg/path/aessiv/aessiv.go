package aessiv

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/google/tink/go/subtle/daead"

	"github.com/connesc/rlink/pkg/path"
)

type Authenticator struct {
	aessiv *daead.AESSIV
}

func NewAuthenticator(key []byte) (*Authenticator, error) {
	aeesiv, err := daead.NewAESSIV(key)
	if err != nil {
		return nil, fmt.Errorf("aessiv: invalid key: %w", err)
	}
	return &Authenticator{
		aessiv: aeesiv,
	}, nil
}

func (r *Authenticator) FromOriginal(originalPath string) (string, error) {
	dir, file := path.Split(originalPath)

	payload := []byte(dir[:len(dir)-1])
	decodedPrefix, err := r.aessiv.EncryptDeterministically(payload, []byte(file))
	if err != nil {
		return "", fmt.Errorf("aessiv: encryption failed: %w", err)
	}

	encodedPrefix := base64.RawURLEncoding.EncodeToString(decodedPrefix)

	return encodedPrefix + "/" + file, nil
}

var authenticatedPathPattern = regexp.MustCompile("^([^/]+)/([^/]*)$")

func (r *Authenticator) ToOriginal(authenticatedPath string) (string, error) {
	matches := authenticatedPathPattern.FindStringSubmatch(authenticatedPath)
	if matches == nil {
		return "", fmt.Errorf("aessiv: invalid authenticated path")
	}
	encodedPrefix := matches[1]
	file := matches[2]

	decodedPrefix := make([]byte, base64.RawURLEncoding.DecodedLen(len(encodedPrefix)))
	_, err := base64.RawURLEncoding.Decode(decodedPrefix, []byte(encodedPrefix))
	if err != nil {
		return "", fmt.Errorf("aessiv: invalid prefix encoding: %w", err)
	}

	payload, err := r.aessiv.DecryptDeterministically(decodedPrefix, []byte(file))
	if err != nil {
		return "", fmt.Errorf("aessiv: invalid prefix encryption: %w", err)
	}

	dir := string(payload) + "/"
	return path.Join(dir, file), nil
}
