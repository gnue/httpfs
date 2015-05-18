package zipfs

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type File struct {
	fi     os.FileInfo
	rc     io.ReadCloser
	files  map[string]*zip.File
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
	finfos, err = f.readdir()
	if err != nil {
		return
	}

	if count < 0 {
		return finfos, nil
	}

	size := int64(len(f.files))

	if size <= f.offset {
		return finfos, io.EOF
	}

	next := f.offset + int64(count)
	if size < next {
		next = size
	}

	finfos = finfos[f.offset:next]
	f.offset += next

	return finfos, nil
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

func (f *File) readdir() ([]os.FileInfo, error) {
	if f.fi.IsDir() && f.files != nil {
		size := int64(len(f.files))

		if f.cache.finfos == nil {
			f.cache.finfos = make([]os.FileInfo, 0, size)

			for _, fname := range f.fnames {
				zf := f.files[fname]
				if zf == nil {
					continue
				}
				f.cache.finfos = append(f.cache.finfos, zf.FileHeader.FileInfo())
			}
		}

		return f.cache.finfos, nil
	}

	return nil, os.ErrNotExist
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

func (f *File) addFile(fn string, zf *zip.File) {
	fname := path.Base(fn)

	if f.files[fname] == nil {
		f.files[fname] = zf
		f.fnames = append(f.fnames, fname)
	}
}
