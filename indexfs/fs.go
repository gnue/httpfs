// Go lang index http.FileSystem
package indexfs

import (
	"net/http"
	"os"
	"path/filepath"
)

type callbackFunc func(http.FileSystem, string) (http.File, error)

type FileSystem struct {
	fileSystem http.FileSystem
	callback   callbackFunc
}

func New(fs http.FileSystem, callback callbackFunc) *FileSystem {
	if callback == nil {
		indexes := Indexes([]string{"index.html"})
		callback = indexes.DirIndex
	}

	return &FileSystem{fileSystem: fs, callback: callback}
}

func (idx *FileSystem) Open(name string) (http.File, error) {
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

func (idx *FileSystem) OpenIndex(dir string) (http.File, error) {
	return idx.callback(idx.fileSystem, dir)
}

func (idx *FileSystem) HasIndex(dir string) bool {
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
