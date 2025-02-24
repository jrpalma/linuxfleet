package html

import (
	"bytes"
	"embed"
	"html/template"
	"log"
)

//go:embed templates
var templateFiles embed.FS

type Templates struct {
	list *template.Template
}

func NewTemplates() *Templates {
	list, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	return &Templates{list}
}

func (t *Templates) Execute(name string, data any) (string, error) {
	output := &bytes.Buffer{}
	err := t.list.ExecuteTemplate(output, name, data)
	return output.String(), err
}
