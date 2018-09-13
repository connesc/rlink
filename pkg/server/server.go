package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/connesc/rlink/pkg/index"
	"github.com/connesc/rlink/pkg/rewriter"
	"github.com/ncw/rclone/vfs"
)

func New(targetPath string, pathRewriter rewriter.PathRewriter) (http.Handler, error) {
	fs, err := NewFs(targetPath)
	if err != nil {
		return nil, err
	}

	handler := &server{
		fs:           fs,
		pathRewriter: pathRewriter,
	}
	return handler, nil
}

type server struct {
	fs           *vfs.VFS
	pathRewriter rewriter.PathRewriter
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	originalPath, err := s.pathRewriter.ToOriginal(req.URL.EscapedPath()[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	node, err := s.fs.Stat(originalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	switch node := node.(type) {
	case *vfs.File:
		if strings.HasSuffix(originalPath, "/") {
			s.redirect(w, req, strings.TrimRight(originalPath, "/"))
			return
		}

		file, err := node.Open(os.O_RDONLY)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
			return
		}

		http.ServeContent(w, req, node.Name(), node.ModTime(), file)

	case *vfs.Dir:
		if !strings.HasSuffix(originalPath, "/") {
			s.redirect(w, req, originalPath+"/")
			return
		}

		children, err := node.ReadDirAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
			return
		}

		content := index.Content{
			Title:   "Index of " + originalPath,
			Entries: make([]index.Entry, len(children)),
		}

		for i, child := range children {
			trailingSlash := ""
			if child.IsDir() {
				trailingSlash = "/"
			}

			nodePath, err := s.pathRewriter.FromOriginal(child.Path() + trailingSlash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
				return
			}

			content.Entries[i] = index.Entry{
				Name: child.Name() + trailingSlash,
				URL:  "/" + nodePath,
			}
		}

		err = index.Write(w, &content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
		}

	default:
		http.Error(w, fmt.Sprintf("unexpected VFS node type: %T", node), http.StatusInternalServerError)
	}
}

func (s *server) redirect(w http.ResponseWriter, req *http.Request, originalPath string) {
	authPath, err := s.pathRewriter.FromOriginal(originalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
		return
	}

	http.Redirect(w, req, "/"+authPath, http.StatusFound)
}
