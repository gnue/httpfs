package unionfs

import (
	"net/http"
	"os"
)

type UnionFS struct {
	fileSystems []http.FileSystem
}

func New(fileSystems ...http.FileSystem) *UnionFS {
	return &UnionFS{fileSystems}
}

func (u *UnionFS) Open(name string) (http.File, error) {
	var dir *Dir

	for _, fs := range u.fileSystems {
		file, err := fs.Open(name)
		if err == nil {
			fi, err := file.Stat()
			if err != nil {
				file.Close()
				return nil, err
			}

			if fi.IsDir() {
				if dir == nil {
					dir = &Dir{fi: &FileInfo{name: fi.Name(), modTime: fi.ModTime()}}
				}

				files, err := file.Readdir(-1)
				file.Close()
				if err != nil {
					continue
				}

				dir.addFile(files...)
			} else {
				if dir != nil {
					file.Close()
					continue
				}

				return file, nil
			}
		}
	}

	if dir != nil {
		return dir, nil
	}

	return nil, os.ErrNotExist
}
