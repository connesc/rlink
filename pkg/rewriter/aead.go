package rewriter

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"strings"
)

type AEADRewriter struct {
	aead      cipher.AEAD
	nonceSize int
	overhead  int
}

func NewAEAD(aead cipher.AEAD) (*AEADRewriter, error) {
	return &AEADRewriter{
		aead:      aead,
		nonceSize: aead.NonceSize(),
		overhead:  aead.Overhead(),
	}, nil
}

func (r *AEADRewriter) FromOriginal(originalPath string) (string, error) {
	trailingSlash := strings.HasSuffix(originalPath, "/")

	dir, file := path.Split(originalPath)
	if len(dir) > 0 {
		dir = dir[:len(dir)-1]
	}

	decodedPrefix := make([]byte, r.nonceSize, r.nonceSize+len(dir)+r.overhead)

	_, err := io.ReadFull(rand.Reader, decodedPrefix)
	if err != nil {
		return "", err
	}

	decodedPrefix = r.aead.Seal(decodedPrefix, decodedPrefix, []byte(dir), []byte(file))
	encodedPrefix := base64.RawURLEncoding.EncodeToString(decodedPrefix)

	return join(encodedPrefix, file, trailingSlash), nil
}

func (r *AEADRewriter) ToOriginal(authPath string) (string, error) {
	trailingSlash := strings.HasSuffix(authPath, "/")

	chunks := strings.SplitN(authPath, "/", 2)

	file := ""
	switch len(chunks) {
	case 0:
		return "", fmt.Errorf("AEADRewriter: no prefix in authenticated path")
	case 2:
		file = chunks[1]
		if strings.Contains(file, "/") {
			return "", fmt.Errorf("AEADRewriter: unexpected slash in file name")
		}
	}
	encodedPrefix := []byte(chunks[0])

	minPrefixSize := r.nonceSize + r.overhead
	prefixSize := base64.RawURLEncoding.DecodedLen(len(encodedPrefix))
	if prefixSize < minPrefixSize {
		return "", fmt.Errorf("AEADRewriter: expected at least %v prefix bytes, go %v", minPrefixSize, prefixSize)
	}

	decodedPrefix := make([]byte, prefixSize)
	_, err := base64.RawURLEncoding.Decode(decodedPrefix, encodedPrefix)
	if err != nil {
		return "", fmt.Errorf("AEADRewriter: invalid prefix encoding: %v", err)
	}

	dir := make([]byte, 0, prefixSize-r.nonceSize-r.overhead)
	dir, err = r.aead.Open(dir, decodedPrefix[:r.nonceSize], decodedPrefix[r.nonceSize:], []byte(file))
	if err != nil {
		return "", fmt.Errorf("AEADRewriter: invalid prefix encryption: %v", err)
	}

	return join(string(dir), file, trailingSlash), nil
}
