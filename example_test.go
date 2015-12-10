package indexfs_test

import (
	"net/http"

	"github.com/gnue/indexfs"
)

func ExampleIndexFS() {
	fs := indexfs.New(http.Dir("sites"), []string{"index.html", "index.htm"})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
