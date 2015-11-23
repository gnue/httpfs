package gitfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type File struct {
	repo   *Repo
	finfo  *FileInfo
	r      *bytes.Reader
	offset int64

	cache struct {
		finfos []os.FileInfo
	}
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(p []byte) (int, error) {
	if f.r == nil {
		return 0, fmt.Errorf("gitfs: can't read %q", f.finfo.Name())
	}

	return f.r.Read(p)
}

func (f *File) Readdir(count int) (finfos []os.FileInfo, err error) {
	finfos, err = f.readdir()
	if err != nil {
		return
	}

	if count < 0 {
		return finfos, nil
	}

	size := int64(len(finfos))

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

func (f *File) readdir() ([]os.FileInfo, error) {
	if !f.finfo.IsDir() {
		return nil, os.ErrNotExist
	}

	if f.cache.finfos == nil {

		args := []string{"ls-tree", "-l", f.finfo.object}

		b, err := f.repo.Exec(args...)
		if err != nil {
			return nil, err
		}

		r := bytes.NewReader(b)
		s := bufio.NewScanner(r)

		finfos := make([]os.FileInfo, 0)

		for s.Scan() {
			finfo, err := parseInfo(s.Text())
			if err != nil {
				return nil, err
			}
			finfos = append(finfos, finfo)
		}

		f.cache.finfos = finfos
	}

	return f.cache.finfos, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.r == nil {
		return 0, fmt.Errorf("gitfs: can't seek %q", f.finfo.Name())
	}

	return f.r.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.finfo, nil
}
