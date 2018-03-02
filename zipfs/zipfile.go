package zipfs

import (
	"archive/zip"
	"net/http"
	"os"
)

type ZipFile struct {
	File *zip.File
}

func (z *ZipFile) FileInfo() os.FileInfo {
	return z.File.FileHeader.FileInfo()
}

func (z *ZipFile) Open() (http.File, error) {
	rc, err := z.File.Open()
	if err != nil {
		return nil, err
	}

	return &File{fi: z.FileInfo(), rc: rc}, nil
}

var _ Finfo = &ZipFile{}
