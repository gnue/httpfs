package indexfs_test

import (
	"net/http"

	"github.com/gnue/httpfs/indexfs"
)

func ExampleIndexFS() {
	fs := indexfs.New(http.Dir("sites"), func(fs http.FileSystem, dir string) (http.File, error) {
		return indexfs.OpenIndex(fs, dir, "index.html", "index.htm")
	})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
