package path

type Path struct {
	authenticator Authenticator
	original      string
	authenticated string
	err           error
}

func New(authenticator Authenticator, original string) *Path {
	return &Path{
		authenticator: authenticator,
		original:      Normalize(original),
	}
}

func NewAuthenticated(authenticator Authenticator, authenticated string) (*Path, error) {
	authenticated = Normalize(authenticated)
	original, err := authenticator.ToOriginal(authenticated)
	if err != nil {
		return nil, err
	}

	path := &Path{
		authenticator: authenticator,
		original:      original,
		authenticated: authenticated,
	}
	return path, nil
}

func (p *Path) Original() string {
	return p.original
}

func (p *Path) Authenticated() (string, error) {
	if p.authenticated == "" && p.err == nil {
		p.authenticated, p.err = p.authenticator.FromOriginal(p.original)
	}
	return p.authenticated, p.err
}

func (p *Path) IsRoot() bool {
	return p.original == "/"
}

func (p *Path) IsDir() bool {
	return IsDir(p.original)
}

func (p *Path) AsDir() *Path {
	original := AsDir(p.original)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}

func (p *Path) AsFile() *Path {
	original := AsFile(p.original)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}

func (p *Path) Dir() *Path {
	original := Dir(p.original)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}

func (p *Path) File() *Path {
	original := File(p.original)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}

func (p *Path) Child(child string) *Path {
	original := Join(p.original, child)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}

func (p *Path) Parent() *Path {
	original := Parent(p.original)
	if original == p.original {
		return p
	}
	return &Path{
		authenticator: p.authenticator,
		original:      original,
	}
}
