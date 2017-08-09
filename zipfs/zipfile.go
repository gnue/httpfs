package zipfs

import (
	"archive/zip"
	"io"
	"os"
)

type ZipFile struct {
	File *zip.File
}

func (z *ZipFile) FileInfo() os.FileInfo {
	return z.File.FileHeader.FileInfo()
}

func (z *ZipFile) Open() (io.ReadCloser, error) {
	return z.File.Open()
}
