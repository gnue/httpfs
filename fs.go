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
	file.Close()

	for _, filename := range idx.indexes {
		fn := filepath.Join(name, filename)

		file, err := fs.Open(fn)
		if err == nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}
