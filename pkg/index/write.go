package index

import "io"

type Entry struct {
	Name string
	URL  string
}

type Content struct {
	Title   string
	Entries []Entry
}

func Write(w io.Writer, content *Content) error {
	return indexTemplate.Execute(w, content)
}
