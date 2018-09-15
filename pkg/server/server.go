package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/connesc/rlink/pkg/index"
	"github.com/connesc/rlink/pkg/path"
	"github.com/ncw/rclone/vfs"
)

type Options struct {
	Files       bool
	Index       bool
	IndexParent bool
}

var defaultOptions = Options{
	Files:       true,
	Index:       true,
	IndexParent: true,
}

func New(targetPath string, authenticator path.Authenticator, options *Options) (http.Handler, error) {
	if options == nil {
		options = &defaultOptions
	}

	fs, err := NewFs(targetPath)
	if err != nil {
		return nil, err
	}

	handler := &server{
		fs:            fs,
		authenticator: authenticator,
		options:       *options,
	}
	return handler, nil
}

type server struct {
	fs            *vfs.VFS
	authenticator path.Authenticator
	options       Options
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath, err := path.NewAuthenticated(s.authenticator, req.URL.EscapedPath())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	node, err := s.fs.Stat(reqPath.Original())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	switch node := node.(type) {
	case *vfs.File:
		if !s.options.Files {
			http.Error(w, "Not Found", http.StatusNotFound) // TODO: better error handling
			return
		}

		if reqPath.IsDir() {
			s.redirect(w, req, reqPath.AsFile())
			return
		}

		file, err := node.Open(os.O_RDONLY)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
			return
		}

		http.ServeContent(w, req, node.Name(), node.ModTime(), file)

	case *vfs.Dir:
		if !s.options.Index {
			http.Error(w, "Not Found", http.StatusNotFound) // TODO: better error handling
			return
		}

		if !reqPath.IsDir() {
			s.redirect(w, req, reqPath.AsDir())
			return
		}

		children, err := node.ReadDirAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
			return
		}

		content := index.Content{
			Title:   "Directory index",
			Entries: make([]index.Entry, 0, 1+len(children)),
		}

		if s.options.IndexParent {
			content.Title = "Index of " + reqPath.Original()
			if !reqPath.IsRoot() {
				parentPath, err := reqPath.Parent().Authenticated()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
					return
				}

				content.Entries = append(content.Entries, index.Entry{
					Name: "../",
					URL:  "/" + parentPath,
				})
			}
		}

		for _, child := range children {
			trailingSlash := ""
			if child.IsDir() {
				trailingSlash = "/"
			}

			childPath, err := path.New(s.authenticator, child.Path()+trailingSlash).Authenticated()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
				return
			}

			content.Entries = append(content.Entries, index.Entry{
				Name: child.Name() + trailingSlash,
				URL:  "/" + childPath,
			})
		}

		err = index.Write(w, &content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
		}

	default:
		http.Error(w, fmt.Sprintf("unexpected VFS node type: %T", node), http.StatusInternalServerError)
	}
}

func (s *server) redirect(w http.ResponseWriter, req *http.Request, dstPath *path.Path) {
	authenticatedPath, err := dstPath.Authenticated()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
		return
	}

	http.Redirect(w, req, path.Absolute(authenticatedPath), http.StatusFound)
}
