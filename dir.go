package unionfs

import (
	"io"
	"os"
	"sort"
	"strings"
)

type Dir struct {
	fi     os.FileInfo
	finfos []os.FileInfo
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

	d.capUp(len(finfos))

	for _, fi := range finfos {
		name := strings.ToLower(fi.Name())
		data := d.finfos

		i := sort.Search(len(data), func(i int) bool {
			return strings.ToLower(data[i].Name()) >= name
		})

		if i < len(data) && strings.ToLower(data[i].Name()) == name {
			continue
		}

		data = append(data, fi)
		copy(data[i+1:], data[i:])
		data[i] = fi

		d.finfos = data
	}
}

func (d *Dir) capUp(n int) {
	l := len(d.finfos)
	if cap(d.finfos) < l+n {
		tmp := make([]os.FileInfo, l, l+n*2)
		copy(tmp, d.finfos)
		d.finfos = tmp
	}
}
