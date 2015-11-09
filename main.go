package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
)

var opts struct {
	Host string `short:"H" long:"host" default:"localhost:3000" description:"host:port"`

	Args struct {
		Dir string `positional-arg-name:"dir" default:"." description:"directory"`
	} `positional-args:"yes"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	fs := http.Dir(opts.Args.Dir)
	err = http.ListenAndServe(opts.Host, http.FileServer(fs))
	if err != nil {
		log.Fatal(err)
	}
}
