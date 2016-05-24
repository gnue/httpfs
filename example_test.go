package indexfs_test

import (
	"net/http"

	"github.com/gnue/httpfs/indexfs"
)

/*
Simple index webserver

```go
indexes := indexfs.Indexes{"index.html", "index.htm"}
fs := indexfs.New(http.Dir("test"), indexes.DirIndex)
```

```go
indexes := indexfs.Indexes{"index.html", "index.htm"}
fs := indexfs.New(http.Dir("test"), indexes.AutoIndex)
```

*/
func ExampleIndexFS() {
	fs := indexfs.New(http.Dir("sites"), func(fs http.FileSystem, dir string) (http.File, error) {
		return indexfs.OpenIndex(fs, dir, "index.html", "index.htm")
	})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
