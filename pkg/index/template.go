package index

import "html/template"

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{ .Title }}</title>
</head>
<body>
	<h1>{{ .Title }}</h1>
	<hr>
	{{ range $i := .Entries }}<a href="{{ $i.URL }}">{{ $i.Name }}</a><br />
	{{ end -}}
	<hr>
</body>
</html>`))
