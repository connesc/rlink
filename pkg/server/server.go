package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/connesc/rlink/pkg/index"
	"github.com/connesc/rlink/pkg/rewriter"
	"github.com/ncw/rclone/vfs"
)

type Options struct {
	Files       bool
	Index       bool
	IndexParent bool
}

var defaultOptions = Options{
	IndexParent: true,
}

func New(targetPath string, pathRewriter rewriter.PathRewriter, options *Options) (http.Handler, error) {
	if options == nil {
		options = &defaultOptions
	}

	fs, err := NewFs(targetPath)
	if err != nil {
		return nil, err
	}

	handler := &server{
		fs:           fs,
		pathRewriter: pathRewriter,
		options:      *options,
	}
	return handler, nil
}

type server struct {
	fs           *vfs.VFS
	pathRewriter rewriter.PathRewriter
	options      Options
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
		if !s.options.Files {
			http.Error(w, "Not Found", http.StatusNotFound) // TODO: better error handling
			return
		}

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
		if !s.options.Index {
			http.Error(w, "Not Found", http.StatusNotFound) // TODO: better error handling
			return
		}

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
			Title:   "Directory index",
			Entries: make([]index.Entry, 0, 1+len(children)),
		}

		if s.options.IndexParent {
			content.Title = "Index of " + originalPath
			if originalPath != "/" {
				parentPath, _ := path.Split(originalPath[:len(originalPath)-1])
				if len(parentPath) == 0 {
					parentPath = "/"
				}

				parentPath, err := s.pathRewriter.FromOriginal(parentPath)
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

			childPath, err := s.pathRewriter.FromOriginal(child.Path() + trailingSlash)
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

func (s *server) redirect(w http.ResponseWriter, req *http.Request, originalPath string) {
	authPath, err := s.pathRewriter.FromOriginal(originalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
		return
	}

	http.Redirect(w, req, "/"+authPath, http.StatusFound)
}
