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
	PageInfo(input []byte) *Page
	Exts() []string
}

type TemplateFS struct {
	Engines      map[string]Engine
	PageTemplete *template.Template
	FileSystem   http.FileSystem
	reExts       *regexp.Regexp
}

func New(fs http.FileSystem, engines ...Engine) *TemplateFS {
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}

	s := strings.TrimLeft(pageTemplate, "\r\n")
	tmpl := template.Must(template.New("generic").Funcs(funcMap).Parse(s))

	t := &TemplateFS{
		FileSystem:   fs,
		Engines:      make(map[string]Engine),
		PageTemplete: tmpl,
	}

	for _, e := range engines {
		t.RegEngine(e)
	}

	t.compileExts()

	return t
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

	for ext, _ := range t.Engines {
		f, err = fs.Open(name + ext)
		if err == nil {
			return f, nil
		}
	}

	return nil, os.ErrNotExist
}

func (t *TemplateFS) RegEngine(e Engine, exts ...string) {
	if exts == nil {
		exts = e.Exts()
	}

	for _, ext := range exts {
		t.Engines[ext] = e
	}
}

func (t *TemplateFS) FindEngine(name string) Engine {
	ext := filepath.Ext(name)
	return t.Engines[ext]
}

func (t *TemplateFS) compileExts() {
	exts := make([]string, 0, len(t.Engines))

	for ext, _ := range t.Engines {
		exts = append(exts, ext)
	}

	t.reExts = regexp.MustCompile(`(?i)\shref="(.+)(` + strings.Join(exts, "|") + `)"`)
}
