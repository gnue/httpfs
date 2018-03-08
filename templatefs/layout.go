package templatefs

import (
	"html/template"
	"strings"
	"unicode"
)

var defaultLayout = template.Must(NewLayout(strings.TrimLeftFunc(pageTemplate, unicode.IsSpace)))

func NewLayout(s string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}

	tmpl := template.New("layout").Funcs(funcMap)
	if s == "" {
		return tmpl, nil
	}

	return tmpl.Parse(s)
}
