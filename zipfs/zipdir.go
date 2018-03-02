package zipfs

import (
	"net/http"
	"os"
)

type ZipDir struct {
	finfo os.FileInfo
}

func (z *ZipDir) FileInfo() os.FileInfo {
	return z.finfo
}

func (z *ZipDir) Open() (http.File, error) {
	return nil, nil
}
