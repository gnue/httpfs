package zipfs

import (
	"archive/zip"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type ZipFS struct {
	Filename string
	rc       *zip.ReadCloser
	dirs     map[string]*File
}

type Options struct {
	Prefix string
	Ignore []string
}

func OpenFS(name string, opts *Options) (z *ZipFS, err error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return
	}

	dirs := make(map[string]*File)

	fi, err := os.Stat(name)
	if err != nil {
		return
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
	dirs["."] = &File{
		fi:    &FileInfo{name: "/", modTime: fi.ModTime()},
		files: make(map[string]Finfo, 0),
	}

	for _, f := range rc.File {
		fi := f.FileHeader.FileInfo()
		fn := strings.Trim(f.FileHeader.Name, "/")
		fn = strings.ToLower(fn)

		// prefix check
		if 0 < len(prefix) {
			if !strings.HasPrefix(fn, prefix) {
				continue
			}
			fn = strings.TrimPrefix(fn, prefix)
			fn = strings.Trim(fn, "/")
		}

		// ignore file
		if ig != nil && ig.MatchString(fn) {
			continue
		}

		if fi.IsDir() {
			if fn == "" {
				fn = "."
			}

			dirs[fn] = &File{
				fi:    fi,
				files: make(map[string]Finfo, 0),
			}

			if fn == "." {
				continue
			}
		}

		dn := path.Dir(fn)
		if dirs[dn] == nil {
			mkpath(dirs, path.Dir(f.FileHeader.Name), fi.ModTime())
		}

		d := dirs[dn]
		d.addFile(fn, &ZipFile{f})
	}

	z = &ZipFS{Filename: name, rc: rc, dirs: dirs}

	return
}

func (z *ZipFS) Open(name string) (file http.File, err error) {
	name = strings.Trim(name, "/")
	if name == "" {
		name = "."
	} else {
		name = strings.ToLower(name)
	}

	d := z.dirs[name]
	if d != nil {
		f := *d
		return &f, nil
	}

	d = z.dirs[path.Dir(name)]
	if d == nil {
		return nil, os.ErrNotExist
	}

	f := d.files[path.Base(name)]
	if f == nil {
		return nil, os.ErrNotExist
	}

	rc, err := f.Open()
	if err != nil {
		return
	}

	return &File{fi: f.FileInfo(), rc: rc}, nil
}

func (z *ZipFS) Close() error {
	return z.rc.Close()
}

func mkpath(dirs map[string]*File, fn string, t time.Time) {
	subdir := strings.Split(fn, "/")
	parent := dirs["."]
	dn := ""

	for _, d := range subdir {
		dn += strings.ToLower(d)

		if dirs[dn] == nil {
			fi := &FileInfo{name: d, modTime: t}
			dirs[dn] = &File{fi: fi, files: make(map[string]Finfo, 0)}
			parent.addFile(dn, &ZipDir{fi})
		}

		parent = dirs[dn]
		dn += "/"
	}
}
