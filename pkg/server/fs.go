package server

import (
	"github.com/ncw/rclone/fs"
	"github.com/ncw/rclone/vfs"

	_ "github.com/ncw/rclone/backend/all" // import all rclone backends
)

func NewFs(targetPath string) (*vfs.VFS, error) {
	backend, err := fs.NewFs(targetPath)
	if err != nil {
		return nil, err
	}
	options := vfs.DefaultOpt
	options.DirCacheTime = 0
	return vfs.New(backend, &options), nil
}
