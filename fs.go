package templatefs

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Engine interface {
	Render(input []byte) []byte
	Exts() []string
}

type TemplateFS struct {
	Engines      []Engine
	PageTemplete *template.Template
	FileSystem   http.FileSystem
	reExts       *regexp.Regexp
}

func New(fs http.FileSystem, e ...Engine) *TemplateFS {
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}

	s := strings.TrimLeft(pageTemplate, "\r\n")
	tmpl := template.Must(template.New("generic").Funcs(funcMap).Parse(s))

	return &TemplateFS{
		FileSystem:   fs,
		Engines:      e,
		PageTemplete: tmpl,
		reExts:       compileExts(e),
	}
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
		file := &File{
			engine:       e,
			pageTemplete: t.PageTemplete,
			file:         f,
			finfo:        finfo,
			reExts:       t.reExts,
		}

		return file, nil
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

func compileExts(engines []Engine) *regexp.Regexp {
	exts := make([]string, 0)

	for _, e := range engines {
		exts = append(exts, e.Exts()...)
	}

	for i, ext := range exts {
		exts[i] = strings.Replace(ext, `.`, `\.`, -1)
	}

	return regexp.MustCompile(`(?i)\shref="(.+)(` + strings.Join(exts, "|") + `)"`)
}
