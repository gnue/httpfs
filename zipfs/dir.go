package zipfs

import (
	"bytes"
	"io"
	"os"
)

type Dir struct {
	*Dinfo
	offset int64

	cache struct {
		finfos []os.FileInfo
		r      *bytes.Reader
	}
}

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Read(p []byte) (int, error) {
	return 0, os.ErrInvalid
}

func (d *Dir) Readdir(count int) (finfos []os.FileInfo, err error) {
	finfos, err = d.readdir()
	if err != nil {
		return
	}

	if count < 0 {
		return finfos, nil
	}

	size := int64(len(d.finfos))

	if size <= d.offset {
		return finfos, io.EOF
	}

	next := d.offset + int64(count)
	if size < next {
		next = size
	}

	finfos = d.finfos[d.offset:next]
	d.offset += next

	return finfos, nil
}

func (d *Dir) readdir() ([]os.FileInfo, error) {
	if d.fi.IsDir() && d.files != nil {
		size := int64(len(d.files))

		if d.cache.finfos == nil {
			d.cache.finfos = make([]os.FileInfo, 0, size)

			for _, fname := range d.fnames {
				fi := d.files[fname]
				if fi == nil {
					continue
				}
				d.cache.finfos = append(d.cache.finfos, fi.FileInfo())
			}
		}

		return d.cache.finfos, nil
	}

	return nil, os.ErrNotExist
}

func (d *Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (d *Dir) Stat() (os.FileInfo, error) {
	return d.fi, nil
}
