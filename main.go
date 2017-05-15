package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gnue/httpfs/gitfs"
	"github.com/gnue/httpfs/indexfs"
	"github.com/gnue/httpfs/unionfs"
	"github.com/gnue/httpfs/zipfs"
	"github.com/jessevdk/go-flags"
	"gopkg.in/gin-gonic/gin.v1"
)

var opts struct {
	Host      string `short:"H" long:"host" default:"localhost" description:"host"`
	Port      string `short:"p" long:"port" default:"3000" description:"port"`
	Branch    string `short:"b" long:"branch" default:"master" description:"git branch"`
	Index     string `long:"index" default:"index.html" description:"directory index"`
	AutoIndex bool   `long:"autoindex" description:"directory auto index"`
	Cert      string `long:"cert" description:"TLS cert file"`
	Key       string `long:"key" description:"TLS key file"`

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

	if !opts.AutoIndex || opts.Index != "index.html" {
		indexes := indexfs.Indexes(strings.Split(opts.Index, ","))
		if opts.AutoIndex {
			fs = indexfs.New(fs, indexes.AutoIndex)
		} else {
			fs = indexfs.New(fs, indexes.DirIndex)
		}
	}

	// Serve
	r := gin.Default()
	r.StaticFS("/", fs)

	if opts.Cert != "" && opts.Key != "" {
		r.RunTLS(addr, opts.Cert, opts.Key)
	} else {
		r.Run(addr)
	}
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

		switch filepath.Ext(d) {
		case ".zip":
			zipOpts := zipfs.Options{Prefix: prefix}
			fs, err := zipfs.OpenFS(d, &zipOpts)
			if err != nil {
				return nil, err
			}
			list[i] = fs
		case ".git":
			list[i] = gitfs.New(d, opts.Branch)
		default:
			list[i] = http.Dir(d)
		}
	}

	return unionfs.New(list...), nil
}
