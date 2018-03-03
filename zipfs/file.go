package zipfs

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	fi     os.FileInfo
	rc     io.ReadCloser
	files  map[string]Finfo
	fnames []string
	offset int64

	cache struct {
		finfos []os.FileInfo
		r      *bytes.Reader
	}
}

func (f *File) Close() error {
	if f.cache.r != nil {
		f.cache.r = nil
	}

	if f.rc == nil {
		return nil
	}

	err := f.rc.Close()
	f.rc = nil

	return err
}

func (f *File) Read(p []byte) (int, error) {
	r, err := f.newReader()
	if err != nil {
		return 0, err
	}

	return r.Read(p)
}

func (f *File) Readdir(count int) (finfos []os.FileInfo, err error) {
	return nil, os.ErrNotExist
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	r, err := f.newReader()
	if err != nil {
		return 0, err
	}

	return r.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

func (f *File) newReader() (*bytes.Reader, error) {
	if f.fi.IsDir() {
		return nil, os.ErrInvalid
	}

	if f.cache.r == nil {
		if f.rc == nil {
			return nil, os.ErrInvalid
		}

		b, err := ioutil.ReadAll(f.rc)
		if err != nil {
			return nil, err
		}

		f.cache.r = bytes.NewReader(b)

		f.rc.Close()
		f.rc = nil
	}

	return f.cache.r, nil
}

func (f *File) addFile(fn string, fi Finfo) {
	fname := filepath.Base(fn)

	if f.files[fname] == nil {
		f.files[fname] = fi
		f.fnames = append(f.fnames, fname)
	}
}
