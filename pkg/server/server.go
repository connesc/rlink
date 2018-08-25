package server

import (
	"net/http"
	"os"

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
	originalPath, err := s.pathRewriter.ToOriginal(req.URL.EscapedPath())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	node, err := s.fs.Stat(originalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}
	if !node.IsFile() {
		http.Error(w, "not a file", http.StatusNotFound) // TODO: support directories
		return
	}

	file, err := node.Open(os.O_RDONLY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // TODO: better error handling
		return
	}

	http.ServeContent(w, req, node.Name(), node.ModTime(), file)
}
