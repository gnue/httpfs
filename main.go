package main

import (
	"github.com/gnue/unionfs"
	"github.com/gnue/zipfs"
	"github.com/jessevdk/go-flags"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var opts struct {
	Host string `short:"H" long:"host" default:"localhost:3000" description:"host:port"`

	Args struct {
		Dir []string `positional-arg-name:"dir" default:"." description:"directory or zip"`
	} `positional-args:"yes"`
}

var reZipPath = regexp.MustCompile(`^(.+\.zip)(?:/(.*))?$`)

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// Address
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

	// File Systems
	dirs := opts.Args.Dir
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	fs, err := newFileSystem(dirs)
	if err != nil {
		log.Fatal(err)
	}

	// Serve
	err = http.ListenAndServe(host, http.FileServer(fs))
	if err != nil {
		log.Fatal(err)
	}
}

func newFileSystem(dirs []string) (http.FileSystem, error) {
	ignore := []string{"__MACOSX", ".DS_Store"}
	list := make([]http.FileSystem, len(dirs))

	for i, d := range dirs {
		prefix := ""

		m := reZipPath.FindStringSubmatch(d)
		if 2 < len(m) {
			d = m[1]
			prefix = m[2]
		}

		if filepath.Ext(d) == ".zip" {
			zipOpts := zipfs.Options{Prefix: prefix, Ignore: ignore}
			fs, err := zipfs.OpenFS(d, &zipOpts)
			if err != nil {
				return nil, err
			}
			list[i] = fs
		} else {
			list[i] = http.Dir(d)
		}
	}

	return unionfs.New(list...), nil
}
