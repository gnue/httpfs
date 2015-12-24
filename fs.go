package indexfs

import (
	"net/http"
	"os"
	"path/filepath"
)

type callbackFunc func(http.FileSystem, string) (http.File, error)

type IndexFS struct {
	fileSystem http.FileSystem
	callback   callbackFunc
}

func New(fs http.FileSystem, callback callbackFunc) *IndexFS {
	if callback == nil {
		callback = func(fs http.FileSystem, dir string) (http.File, error) {
			return OpenIndex(fs, dir, "index.html")
		}
	}

	return &IndexFS{fileSystem: fs, callback: callback}
}

func (idx *IndexFS) Open(name string) (http.File, error) {
	fs := idx.fileSystem

	file, err := fs.Open(name)
	if err != nil {
		if name == "/index.html" {
			return idx.OpenIndex("/")
		}

		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	if !fi.IsDir() {
		return file, nil
	}
	if name == "/" {
		if idx.HasIndex(name) {
			return file, nil
		}
		file.Close()
		return nil, os.ErrNotExist
	}
	file.Close()

	return idx.OpenIndex(name)
}

func (idx *IndexFS) OpenIndex(dir string) (http.File, error) {
	return idx.callback(idx.fileSystem, dir)
}

func (idx *IndexFS) HasIndex(dir string) bool {
	f, err := idx.OpenIndex(dir)
	if err == nil {
		f.Close()
		return true
	}

	return false
}

func OpenIndex(fs http.FileSystem, dir string, indexes ...string) (http.File, error) {
	for _, filename := range indexes {
		fn := filepath.Join(dir, filename)

		file, err := fs.Open(fn)
		if err == nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}
