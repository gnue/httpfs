package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gnue/indexfs"
	"github.com/gnue/unionfs"
	"github.com/gnue/zipfs"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Host  string `short:"H" long:"host" default:"localhost" description:"host"`
	Port  string `short:"p" long:"port" default:"3000" description:"port"`
	Index string `long:"index" default:"index.html" description:"directory index"`

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
	addr := net.JoinHostPort(opts.Host, opts.Port)

	// File Systems
	dirs := opts.Args.Dir
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	fs, err := newFileSystem(dirs)
	if err != nil {
		log.Fatal(err)
	}

	switch opts.Index {
	case "false":
	case "true", "":
		fs = indexfs.New(fs, nil)
	default:
		fs = indexfs.New(fs, strings.Split(opts.Index, ","))
	}

	// Serve
	r := gin.Default()
	r.StaticFS("/", fs)
	r.Run(addr)
}

func newFileSystem(dirs []string) (http.FileSystem, error) {
	list := make([]http.FileSystem, len(dirs))

	for i, d := range dirs {
		prefix := ""

		m := reZipPath.FindStringSubmatch(d)
		if 2 < len(m) {
			d = m[1]
			prefix = m[2]
		}

		if filepath.Ext(d) == ".zip" {
			zipOpts := zipfs.Options{Prefix: prefix}
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
