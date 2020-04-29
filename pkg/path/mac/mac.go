package mac

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"regexp"

	"golang.org/x/crypto/sha3"
)

type MAC interface {
	Compute(payload []byte) []byte
	Size() int
}

type MACFunc func(key []byte) (MAC, error)

type hashMAC struct {
	hash.Hash
}

func (m hashMAC) Compute(payload []byte) []byte {
	m.Write(payload)
	return m.Sum(nil)
}

func HMAC(hashFunc func() hash.Hash) MACFunc {
	return func(key []byte) (MAC, error) {
		return hashMAC{hmac.New(hashFunc, key)}, nil
	}
}

func KeyedHash(hashFunc func(key []byte) (hash.Hash, error)) MACFunc {
	return func(key []byte) (MAC, error) {
		hash, err := hashFunc(key)
		if err != nil {
			return nil, err
		}
		return hashMAC{hash}, nil
	}
}

func VariableKeyedHash(hashFunc func(size int, key []byte) (hash.Hash, error), size int) MACFunc {
	return func(key []byte) (MAC, error) {
		hash, err := hashFunc(size, key)
		if err != nil {
			return nil, err
		}
		return hashMAC{hash}, nil
	}
}

func PrefixedHash(hashFunc func() hash.Hash) MACFunc {
	return func(key []byte) (MAC, error) {
		hash := hashFunc()
		hash.Write(key)
		return hashMAC{hash}, nil
	}
}

type shakeHashMAC struct {
	sha3.ShakeHash
	size int
}

func (m shakeHashMAC) Compute(payload []byte) []byte {
	m.Write(payload)
	mac := make([]byte, m.size)
	m.Read(mac)
	return mac
}

func (m shakeHashMAC) Size() int {
	return m.size
}

func PrefixedShakeHash(hashFunc func() sha3.ShakeHash, size int) MACFunc {
	return func(key []byte) (MAC, error) {
		hash := hashFunc()
		hash.Write(key)
		return shakeHashMAC{hash, size}, nil
	}
}

type Authenticator struct {
	macFunc MACFunc
	key     []byte
}

func NewAuthenticator(macFunc MACFunc, key []byte) *Authenticator {
	return &Authenticator{
		macFunc: macFunc,
		key:     key,
	}
}

func (r *Authenticator) FromOriginal(originalPath string) (string, error) {
	payload := originalPath
	if payload == "/" {
		payload = ""
	}

	mac, err := r.macFunc(r.key)
	if err != nil {
		return "", fmt.Errorf("mac: %w", err)
	}

	decodedMAC := mac.Compute([]byte(payload))
	encodedMAC := base64.RawURLEncoding.EncodeToString(decodedMAC)

	return encodedMAC + "/" + payload, nil
}

var authenticatedPathPattern = regexp.MustCompile("^([^/]+)/(.*)$")

func (r *Authenticator) ToOriginal(authenticatedPath string) (string, error) {
	matches := authenticatedPathPattern.FindStringSubmatch(authenticatedPath)
	if matches == nil {
		return "", fmt.Errorf("mac: invalid authenticated path")
	}
	encodedMAC := matches[1]
	payload := matches[2]

	mac, err := r.macFunc(r.key)
	if err != nil {
		return "", fmt.Errorf("mac: %w", err)
	}

	expectedMACBytes := mac.Size()

	actualMACBytes := base64.RawURLEncoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("mac: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err = base64.RawURLEncoding.Decode(decodedMAC, []byte(encodedMAC))
	if err != nil {
		return "", fmt.Errorf("mac: invalid MAC: %w", err)
	}

	computedMAC := mac.Compute([]byte(payload))

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("mac: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	originalPath := payload
	if originalPath == "" {
		originalPath = "/"
	}
	return originalPath, nil
}
