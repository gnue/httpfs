package gitfs_test

import (
	"net/http"

	"github.com/gnue/httpfs/gitfs"
)

// Simple git webserver
func ExampleGitFS() {
	fs := gitfs.New("sites.gif", "master")

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
