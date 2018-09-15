package path

import (
	"regexp"
	"strings"
)

var repeatedSlashes = regexp.MustCompile("/+")

func Normalize(path string) string {
	path = repeatedSlashes.ReplaceAllLiteralString(path, "/")
	if path == "" || path == "/" {
		return "/"
	}
	return strings.TrimPrefix(path, "/")
}

func Absolute(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func IsDir(path string) bool {
	return strings.HasSuffix(path, "/")
}

func AsDir(path string) string {
	if IsDir(path) {
		return path
	}
	return path + "/"
}

func AsFile(path string) string {
	if path == "/" {
		return path
	}
	return strings.TrimSuffix(path, "/")
}

func Split(path string) (string, string) {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return "/", path
	}
	return path[:lastSlash+1], path[lastSlash+1:]
}

func Dir(path string) string {
	dir, _ := Split(path)
	return dir
}

func File(path string) string {
	_, file := Split(path)
	return file
}

func Join(parent, child string) string {
	parent = AsDir(parent)
	if child == "" || child == "/" {
		return parent
	}
	if parent == "/" {
		return child
	}
	return parent + child
}

func Parent(path string) string {
	return Dir(AsFile(path))
}
