package path

type Authenticator interface {
	FromOriginal(path string) (string, error)
	ToOriginal(path string) (string, error)
}
