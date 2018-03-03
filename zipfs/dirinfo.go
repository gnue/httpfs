package zipfs

import (
	"os"
	"time"
)

type DirInfo struct {
	name    string
	modTime time.Time
}

func (f *DirInfo) Name() string {
	return f.name
}

func (f *DirInfo) Size() int64 {
	return 0
}

func (f *DirInfo) Mode() os.FileMode {
	return os.ModeDir | 0755
}

func (f *DirInfo) IsDir() bool {
	return true
}

func (f *DirInfo) ModTime() time.Time {
	return f.modTime
}

func (f *DirInfo) Sys() interface{} {
	return nil
}
