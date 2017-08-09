# unionfs

Go lang union http.FileSystem

## Installation

```sh
$ go get github.com/gnue/httpfs/unionfs
```

## Usage

```go
import "github.com/gnue/httpfs/unionfs"
```

## Examples

### UnionFS

```go
package main

import (
	"github.com/gnue/httpfs/unionfs"
	"net/http"
)

func main() {
	fs := unionfs.New(http.Dir("A"), http.Dir("B"))

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple union file system webserver

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

