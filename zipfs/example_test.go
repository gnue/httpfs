package zipfs_test

import (
	"log"
	"net/http"

	"github.com/gnue/httpfs/zipfs"
)

// Simple zip webserver
func ExampleOpen() {
	fs, err := zipfs.Open("public.zip", &zipfs.Options{Prefix: "public"})
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
