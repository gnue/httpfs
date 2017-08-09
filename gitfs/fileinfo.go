package gitfs

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileInfo struct {
	mode    string
	objtype string
	object  string
	size    int64
	path    string

	modTime time.Time
}

func (f *FileInfo) Name() string {
	return filepath.Base(f.path)
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() os.FileMode {
	i, _ := strconv.Atoi(f.mode)
	m := os.FileMode(i & 0777)

	if f.objtype == "tree" {
		m &= os.ModeDir
	}

	return m
}

func (f *FileInfo) IsDir() bool {
	if f.objtype == "tree" {
		return true
	}

	return false
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) Sys() interface{} {
	return nil
}
