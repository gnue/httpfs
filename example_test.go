package unionfs_test

import (
	"net/http"

	"github.com/gnue/httpfs/unionfs"
)

// Simple union file system webserver
func ExampleUnionFS() {
	fs := unionfs.New(http.Dir("A"), http.Dir("B"))

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
