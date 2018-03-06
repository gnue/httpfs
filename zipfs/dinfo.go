package zipfs

import (
	"net/http"
	"os"
	"path/filepath"
)

type Dinfo struct {
	fi     os.FileInfo
	finfos []os.FileInfo
	files  map[string]Finfo
	fnames []string
}

func newDir(fi os.FileInfo) *Dinfo {
	return &Dinfo{fi: fi, files: make(map[string]Finfo, 0)}
}

func (d *Dinfo) FileInfo() os.FileInfo {
	return d.fi
}

func (d *Dinfo) Open() (http.File, error) {
	return &Dir{Dinfo: d}, nil
}

func (d *Dinfo) addFile(fn string, fi Finfo) {
	fname := filepath.Base(fn)

	if d.files[fname] == nil {
		d.files[fname] = fi
		d.fnames = append(d.fnames, fname)
	}
}

var _ Finfo = &Dinfo{}
