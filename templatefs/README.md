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
	fs := templatefs.New(http.Dir("templates"), &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple markdown webserver

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

