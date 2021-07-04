# templatefs

Go lang template http.FileSystem

## Installation

```sh
$ go get github.com/gnue/httpfs/templatefs
```

## Usage

```go
import "github.com/gnue/httpfs/templatefs"
```

## Examples

### TemplateFS

```go
package main

import (
	"github.com/gnue/httpfs/templatefs"
	"github.com/gnue/httpfs/templatefs/engines/markdown"
	"net/http"
)

func main() {
	fs := templatefs.New(http.Dir("templates"), nil, &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple markdown webserver

### TemplateFS_custom_layout

```go
package main

import (
	"github.com/gnue/httpfs/templatefs"
	"github.com/gnue/httpfs/templatefs/engines/markdown"
	"net/http"
	"strings"
	"unicode"
)

func main() {
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

	layout, err := templatefs.NewLayout(strings.TrimLeftFunc(s, unicode.IsSpace))
	if err != nil {
		return
	}

	fs := templatefs.New(http.Dir("templates"), layout, &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple markdown webserver with custom layout

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

