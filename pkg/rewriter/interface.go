package rewriter

type PathRewriter interface {
	FromOriginal(path string) (string, error)
	ToOriginal(path string) (string, error)
}
