package templatefs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
)

type renderFunc func([]byte) ([]byte, error)

type File struct {
	file   http.File
	finfo  os.FileInfo
	r      *bytes.Reader
	render renderFunc
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Read(p []byte) (int, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return 0, err
		}
		f.r = r
	}

	i, err := f.r.Read(p)
	return i, err

	return f.r.Read(p)
}

func (f *File) newReader() (*bytes.Reader, error) {
	if f.finfo.IsDir() {
		return nil, os.ErrInvalid
	}

	b, err := ioutil.ReadAll(f.file)
	if err != nil {
		return nil, err
	}
	f.file.Seek(0, os.SEEK_CUR)

	output, err := f.render(b)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(output), nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return 0, err
		}
		f.r = r
	}

	return f.r.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return nil, err
		}
		f.r = r
	}

	return &FileInfo{finfo: f.finfo, size: f.r.Size()}, nil
}
