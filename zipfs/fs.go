// Go lang zip http.FileSystem
package zipfs

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var OSX_Ignore = []string{"__MACOSX", ".DS_Store"}

type FileSystem struct {
	r    Reader
	dirs map[string]*Dinfo
}

type FileSystemCloser struct {
	FileSystem
	Filename string
	c        io.Closer
}

type Options struct {
	Prefix string
	Ignore []string
}

type Finfo interface {
	FileInfo() os.FileInfo
	Open() (http.File, error)
}

func New(data []byte, opts *Options) (*FileSystem, error) {
	b := bytes.NewReader(data)
	r, err := zip.NewReader(b, b.Size())
	if err != nil {
		return nil, err
	}

	dirs := newDirs(r.File, time.Now(), opts)

	return &FileSystem{&reader{r}, dirs}, nil
}

func Open(name string, opts *Options) (*FileSystemCloser, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(name)
	if err != nil {
		rc.Close()
		return nil, err
	}

	dirs := newDirs(rc.File, fi.ModTime(), opts)
	fs := FileSystem{&readCloser{rc}, dirs}

	return &FileSystemCloser{fs, name, rc}, nil
}

func (fs *FileSystem) Open(name string) (file http.File, err error) {
	name = strings.Trim(name, "/")
	if name == "" {
		name = "."
	} else {
		name = strings.ToLower(name)
	}

	d := fs.dirs[name]
	if d != nil {
		return d.Open()
	}

	d = fs.dirs[filepath.Dir(name)]
	if d == nil {
		return nil, os.ErrNotExist
	}

	f := d.files[filepath.Base(name)]
	if f == nil {
		return nil, os.ErrNotExist
	}

	return f.Open()
}

func (z *FileSystemCloser) Close() error {
	return z.c.Close()
}

func newDirs(files []*zip.File, modTime time.Time, opts *Options) map[string]*Dinfo {
	dirs := make(map[string]*Dinfo)

	// opts
	if opts == nil {
		opts = &Options{}
	}

	if opts.Ignore == nil {
		opts.Ignore = OSX_Ignore
	}

	// ignore files
	ig, _ := NewIgnore(opts.Ignore)

	// prefix
	prefix := opts.Prefix
	if 0 < len(prefix) {
		prefix = strings.ToLower(prefix)
		prefix = strings.Trim(prefix, "/") + "/"
	}

	// root directory
	dirs["."] = newDir(&DirInfo{name: "/", modTime: modTime})

	for _, f := range files {
		fi := f.FileHeader.FileInfo()
		org := strings.Trim(f.FileHeader.Name, "/")
		fn := strings.ToLower(org)

		// prefix check
		if 0 < len(prefix) {
			if !strings.HasPrefix(fn, prefix) {
				continue
			}
			fn = strings.TrimPrefix(fn, prefix)
			fn = strings.Trim(fn, "/")
			org = org[len(org)-len(fn):]
		}

		// ignore file
		if ig != nil && ig.MatchString(fn) {
			continue
		}

		if fi.IsDir() {
			if fn == "" {
				fn = "."
			}

			dirs[fn] = newDir(fi)

			if fn == "." {
				continue
			}
		}

		dn := filepath.Dir(fn)
		if dirs[dn] == nil {
			mkpath(dirs, filepath.Dir(org), fi.ModTime())
		}

		d := dirs[dn]
		d.addFile(fn, &ZipFile{f})
	}

	return dirs
}

func mkpath(dirs map[string]*Dinfo, fn string, t time.Time) {
	subdir := strings.Split(fn, "/")
	parent := dirs["."]
	dn := ""

	for _, d := range subdir {
		dn += strings.ToLower(d)

		if dirs[dn] == nil {
			fi := &DirInfo{name: d, modTime: t}
			dirs[dn] = newDir(fi)
			if parent != nil {
				parent.addFile(dn, &ZipDir{fi})
			}
		}

		parent = dirs[dn]
		dn += "/"
	}
}
