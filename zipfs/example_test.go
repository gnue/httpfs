package zipfs_test

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gnue/httpfs/zipfs"
)

// Simple zip webserver(use New)
func ExampleNew() {
	b, err := ioutil.ReadFile("public.zip")
	if err != nil {
		log.Fatal(err)
	}

	fs, err := zipfs.New(b, &zipfs.Options{Prefix: "public"})
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}

// Simple zip webserver(use Open)
func ExampleOpen() {
	fs, err := zipfs.Open("public.zip", &zipfs.Options{Prefix: "public"})
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	http.Handle("/", http.FileServer(fs))
	http.ListenAndServe(":8080", nil)
}
