package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"strings"
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

	h := strings.Split(opts.Host, ":")
	if len(h) < 2 {
		h = append(h, "")
	}

	if h[0] == "" {
		h[0] = "localhost"
	}
	if h[1] == "" {
		h[1] = "80"
	}

	host := strings.Join(h, ":")

	fs := http.Dir(opts.Args.Dir)
	err = http.ListenAndServe(host, http.FileServer(fs))
	if err != nil {
		log.Fatal(err)
	}
}
