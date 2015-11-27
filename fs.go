package templatefs

import (
	"net/http"
	"os"
	"path/filepath"
)

type Engine interface {
	Render(input []byte) []byte
	Exts() []string
}

type TemplateFS struct {
	Engines    []Engine
	FileSystem http.FileSystem
}

func New(fs http.FileSystem, e ...Engine) *TemplateFS {
	return &TemplateFS{FileSystem: fs, Engines: e}
}

func (t *TemplateFS) Open(name string) (http.File, error) {
	f, err := t.open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if finfo.IsDir() {
		return f, nil
	}

	e := t.FindEngine(finfo.Name())
	if e != nil {
		return &File{engine: e, file: f, finfo: finfo}, nil
	}

	return f, nil
}

func (t *TemplateFS) open(name string) (http.File, error) {
	fs := t.FileSystem
	f, err := fs.Open(name)
	if err == nil {
		return f, nil
	}

	for _, e := range t.Engines {
		for _, ext := range e.Exts() {
			f, err = fs.Open(name + ext)
			if err == nil {
				return f, nil
			}
		}
	}

	return nil, os.ErrNotExist
}

func (t *TemplateFS) FindEngine(name string) Engine {
	ext := filepath.Ext(name)

	for _, e := range t.Engines {
		for _, v := range e.Exts() {
			if ext == v {
				return e
			}
		}
	}

	return nil
}
