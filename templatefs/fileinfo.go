package templatefs

import (
	"os"
	"time"
)

type FileInfo struct {
	finfo os.FileInfo
	size  int64
	name  string

	modTime time.Time
}

func (f *FileInfo) Name() string {
	if f.name != "" {
		return f.name
	}

	return f.finfo.Name()
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() os.FileMode {
	return f.finfo.Mode()
}

func (f *FileInfo) IsDir() bool {
	return false
}

func (f *FileInfo) ModTime() time.Time {
	return f.finfo.ModTime()
}

func (f *FileInfo) Sys() interface{} {
	return f.finfo.Sys()
}
