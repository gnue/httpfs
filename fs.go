package indexfs

import (
	"net/http"
	"os"
	"path/filepath"
)

var DefaultIndexes = []string{"index.html"}

type IndexFS struct {
	fileSystem http.FileSystem
	indexes    []string
}

func New(fs http.FileSystem, indexes []string) *IndexFS {
	if indexes == nil {
		indexes = DefaultIndexes
	}

	return &IndexFS{fileSystem: fs, indexes: indexes}
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
	fs := idx.fileSystem

	for _, filename := range idx.indexes {
		fn := filepath.Join(dir, filename)

		file, err := fs.Open(fn)
		if err == nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}

func (idx *IndexFS) HasIndex(dir string) bool {
	f, err := idx.OpenIndex(dir)
	if err == nil {
		f.Close()
		return true
	}

	return false
}
