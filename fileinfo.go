package zipfs

import (
	"os"
	"time"
)

type FileInfo struct {
	name    string
	modTime time.Time
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return 0
}

func (f *FileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (f *FileInfo) IsDir() bool {
	return true
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) Sys() interface{} {
	return nil
}
