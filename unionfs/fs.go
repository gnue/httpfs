// Go lang union http.FileSystem
package unionfs

import (
	"net/http"
	"os"
)

type FileSystem struct {
	fileSystems []http.FileSystem
}

func New(fileSystems ...http.FileSystem) *FileSystem {
	return &FileSystem{fileSystems}
}

func (u *FileSystem) Open(name string) (http.File, error) {
	var dirFSs []http.FileSystem

	for _, fs := range u.fileSystems {
		file, err := fs.Open(name)
		if err != nil {
			continue
		}

		fi, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, err
		}

		if !fi.IsDir() && dirFSs == nil {
			return file, nil
		}

		file.Close()
		dirFSs = append(dirFSs, fs)
	}

	if len(dirFSs) == 1 {
		return dirFSs[0].Open(name)
	}

	if dirFSs != nil {
		var dir *Dir

		for _, fs := range dirFSs {
			file, err := fs.Open(name)
			if err != nil {
				continue
			}

			var files []os.FileInfo

			fi, err := file.Stat()
			if err == nil {
				files, err = file.Readdir(-1)
			}
			file.Close()

			if err != nil {
				continue
			}

			if dir == nil {
				dir = &Dir{fi: &DirInfo{name: fi.Name(), modTime: fi.ModTime()}}
			}

			dir.addFile(files...)
		}

		if dir != nil {
			return dir, nil
		}
	}

	return nil, os.ErrNotExist
}
