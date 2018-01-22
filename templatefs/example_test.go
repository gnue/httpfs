package templatefs_test

import (
	"net/http"
	"strings"

	"github.com/gnue/httpfs/templatefs"
	"github.com/gnue/httpfs/templatefs/engines/markdown"
)

// Simple markdown webserver
func ExampleTemplateFS() {
	fs := templatefs.New(http.Dir("templates"), nil, &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

// Simple markdown webserver with custom layout
func ExampleTemplateFS_custom_layout() {
	s := `
<!DOCTYPE html>
<html>
<head>
	<title>{{ .Title }}</title>
	<meta charset="utf-8">
	{{- if .CSS}}
	<link rel="stylesheet" type="text/css" href="{{ .CSS }}">
	{{- end}}
</head>
<body>
{{ .Body | safehtml }}
</body>
</html>
`

	layout, err := templatefs.NewLayout(strings.TrimLeft(s, "\r\n"))
	if err != nil {
		return
	}

	fs := templatefs.New(http.Dir("templates"), layout, &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
