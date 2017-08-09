# indexfs

Go lang index http.FileSystem

## Installation

```sh
$ go get github.com/gnue/httpfs/indexfs
```

## Usage

```go
import "github.com/gnue/httpfs/indexfs"
```

## Examples

### IndexFS

```go
package main

import (
	"github.com/gnue/httpfs/indexfs"
	"net/http"
)

func main() {
	fs := indexfs.New(http.Dir("sites"), func(fs http.FileSystem, dir string) (http.File, error) {
		return indexfs.OpenIndex(fs, dir, "index.html", "index.htm")
	})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple index webserver

```go
indexes := indexfs.Indexes{"index.html", "index.htm"}
fs := indexfs.New(http.Dir("test"), indexes.DirIndex)
```

```go
indexes := indexfs.Indexes{"index.html", "index.htm"}
fs := indexfs.New(http.Dir("test"), indexes.AutoIndex)
```

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

