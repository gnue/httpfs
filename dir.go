package unionfs

import (
	"io"
	"os"
)

type Dir struct {
	fi     os.FileInfo
	finfos []os.FileInfo
	files  map[string]os.FileInfo
	offset int64
}

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Read(p []byte) (int, error) {
	return 0, os.ErrInvalid
}

func (d *Dir) Readdir(count int) (finfos []os.FileInfo, err error) {
	finfos = d.finfos

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

func (d *Dir) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}

func (d *Dir) Stat() (os.FileInfo, error) {
	return d.fi, nil
}

func (d *Dir) addFile(finfos ...os.FileInfo) {
	if len(finfos) == 0 {
		return
	}

	if d.files == nil {
		d.files = make(map[string]os.FileInfo)
	}

	for _, fi := range finfos {
		name := fi.Name()
		if d.files[name] != nil {
			continue
		}

		d.files[name] = fi
		d.finfos = append(d.finfos, fi)
	}
}
