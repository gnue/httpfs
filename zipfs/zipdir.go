package zipfs

import (
	"io"
	"os"
)

type ZipDir struct {
	finfo os.FileInfo
}

func (z *ZipDir) FileInfo() os.FileInfo {
	return z.finfo
}

func (z *ZipDir) Open() (io.ReadCloser, error) {
	return nil, nil
}
