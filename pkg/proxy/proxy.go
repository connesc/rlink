package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/connesc/rlink/pkg/path"
)

func New(targetURL string, authenticator path.Authenticator) (http.Handler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	targetQuery := target.RawQuery
	targetPath := path.Normalize(target.EscapedPath())

	directorWithErr := func(req *http.Request) (err error) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		originalPath, err := authenticator.ToOriginal(path.Normalize(req.URL.EscapedPath()))
		if err != nil {
			return
		}

		req.URL.RawPath = path.Absolute(path.Join(targetPath, originalPath))
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

	return &httputil.ReverseProxy{Director: director}, nil
}
