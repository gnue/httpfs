# zipfs

Go lang zip http.FileSystem

## Installation

```sh
$ go get github.com/gnue/httpfs/zipfs
```

## Usage

```go
import "github.com/gnue/httpfs/zipfs"
```

## Examples

### ZipFS

```go
package main

import (
	"github.com/gnue/httpfs/zipfs"
	"log"
	"net/http"
)

func main() {
	fs, err := zipfs.OpenFS("public.zip", &zipfs.Options{Prefix: "public"})
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

```

Simple zip webserver

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

