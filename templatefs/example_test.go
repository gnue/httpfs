package templatefs_test

import (
	"net/http"

	"github.com/gnue/httpfs/templatefs"
	"github.com/gnue/httpfs/templatefs/engines/markdown"
)

// Simple markdown webserver
func ExampleTemplateFS() {
	fs := templatefs.New(http.Dir("templates"), &markdown.Engine{})

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
