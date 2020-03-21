package server

import (
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/vfs"

	_ "github.com/rclone/rclone/backend/all" // import all rclone backends
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
