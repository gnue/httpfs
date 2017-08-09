package fsutil

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var SkipDir = filepath.SkipDir

type FileSystem struct {
	http.FileSystem
}

func (fs *FileSystem) Stat(name string) (os.FileInfo, error) {
	type call interface {
		Stat(string) (os.FileInfo, error)
	}

	if fs, ok := fs.FileSystem.(call); ok {
		return fs.Stat(name)
	}

	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Stat()
}

func (fs *FileSystem) Glob(pattern string) (matches []string, err error) {
	if !hasMeta(pattern) {
		if _, err := fs.Stat(pattern); err == nil {
			return nil, nil
		}
		return []string{pattern}, nil
	}

	dir, file := filepath.Split(pattern)
	switch dir {
	case "":
		dir = "."
	case string(filepath.Separator):
		// nothing
	default:
		dir = dir[:len(dir)-1]
	}

	if !hasMeta(dir) {
		return fs.glob(dir, file, nil)
	}

	m, err := fs.Glob(dir)
	if err != nil {
		return
	}

	for _, d := range m {
		matches, err = fs.glob(d, file, matches)
		if err != nil {
			return
		}
	}

	sort.Strings(matches)

	return
}

func (fs *FileSystem) glob(dir, pattern string, matches []string) (m []string, err error) {
	m = matches

	finfo, err := fs.Stat(dir)
	if err != nil {
		return
	}
	if !finfo.IsDir() {
		return
	}

	if pattern == "**" {
		err = fs.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if dir != path {
				m = append(m, path)
			}
			return nil
		})

		return
	}

	names, err := fs.Readdirnames(dir)
	if err != nil {
		return
	}

	for _, fname := range names {
		matched, err := filepath.Match(pattern, fname)
		if err != nil {
			return m, err
		}
		if matched {
			m = append(m, filepath.Join(dir, fname))
		}
	}

	return
}

func hasMeta(path string) bool {
	return 0 <= strings.IndexAny(path, "*?[")
}

func (fs *FileSystem) Walk(root string, walkFn filepath.WalkFunc) error {
	finfo, err := fs.Stat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}

	return fs.walk(root, finfo, walkFn)
}

func (fs *FileSystem) walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := fs.Readdirnames(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		fname := filepath.Join(path, name)
		finfo, err := fs.Stat(fname)
		if err != nil {
			err = walkFn(fname, finfo, err)
			if err != nil && err != SkipDir {
				return err
			}
		} else {
			err = fs.walk(fname, finfo, walkFn)
			if err != nil {
				if !finfo.IsDir() || err != SkipDir {
					return err
				}
			}
		}
	}

	return nil
}

func (fs *FileSystem) Readdirnames(name string) ([]string, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	infos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(infos))
	for i, finfo := range infos {
		names[i] = finfo.Name()
	}
	sort.Strings(names)

	return names, nil
}

func (fs *FileSystem) Readfile(name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}
