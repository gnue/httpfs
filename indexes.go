package indexfs

import (
	"net/http"
)

type Indexes []string

func (idx Indexes) DirIndex(fs http.FileSystem, dir string) (http.File, error) {
	return OpenIndex(fs, dir, idx...)
}

func (idx Indexes) AutoIndex(fs http.FileSystem, dir string) (http.File, error) {
	f, err := OpenIndex(fs, dir, idx...)
	if err == nil {
		return f, nil
	}

	return fs.Open(dir)
}
