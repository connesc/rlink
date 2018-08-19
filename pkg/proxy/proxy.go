package proxy

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func New(targetURL string, secret []byte) http.Handler {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	targetQuery := target.RawQuery
	targetPath := target.EscapedPath()

	if !strings.HasSuffix(targetPath, "/") {
		targetPath += "/"
	}

	directorWithErr := func(req *http.Request) (err error) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		req.URL.RawPath, err = rewritePath(targetPath, req.URL.EscapedPath(), sha1.New, secret, base64.RawURLEncoding)
		if err != nil {
			return
		}

		req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
		if err != nil {
			return
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("Host", target.Host)

		return
	}

	director := func(req *http.Request) {
		err := directorWithErr(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			req.URL = nil
			return
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

func rewritePath(basePath string, reqPath string, hashFunc func() hash.Hash, secret []byte, encoding *base64.Encoding) (string, error) {
	chunks := strings.SplitN(reqPath[1:], "/", 2)

	newPath := ""
	switch len(chunks) {
	case 0:
		return "", fmt.Errorf("rewritePath: no MAC in request URL")
	case 2:
		newPath = chunks[1]
	}
	encodedMAC := []byte(chunks[0])

	mac := hmac.New(hashFunc, secret)
	expectedMACBytes := mac.Size()

	actualMACBytes := encoding.DecodedLen(len(encodedMAC))
	if actualMACBytes != expectedMACBytes {
		return "", fmt.Errorf("rewritePath: invalid MAC: expected %v bytes, got %v", expectedMACBytes, actualMACBytes)
	}

	decodedMAC := make([]byte, expectedMACBytes)
	_, err := encoding.Decode(decodedMAC, encodedMAC)
	if err != nil {
		return "", fmt.Errorf("rewritePath: invalid MAC: %v", err)
	}

	mac.Write([]byte(newPath))
	computedMAC := mac.Sum(nil)

	if !hmac.Equal(decodedMAC, computedMAC) {
		return "", fmt.Errorf("rewritePath: MAC does not match (%v != %v)", decodedMAC, computedMAC)
	}

	return basePath + newPath, nil
}
