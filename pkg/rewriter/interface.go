package rewriter

type URLRewriter interface {
	FromOriginal(path string) (string, error)
	ToOriginal(path string) (string, error)
}
