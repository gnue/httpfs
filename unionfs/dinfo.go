package unionfs

import (
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Dinfo struct {
	fi     os.FileInfo
	finfos []os.FileInfo
}

func newDir(name string, modTime time.Time) *Dinfo {
	return &Dinfo{fi: &DirInfo{name: name, modTime: modTime}}
}

func (d *Dinfo) Open() (http.File, error) {
	return &Dir{Dinfo: d}, nil
}

func (d *Dinfo) addFile(finfos ...os.FileInfo) {
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

func (d *Dinfo) capUp(n int) {
	l := len(d.finfos)
	if cap(d.finfos) < l+n {
		tmp := make([]os.FileInfo, l, l+n*2)
		copy(tmp, d.finfos)
		d.finfos = tmp
	}
}
