package html

import (
	"bytes"
	"html/template"
	"log"
)

type Templates struct {
	list *template.Template
}

func NewTemplates() *Templates {
	list, err := template.ParseGlob("templates/*.tmpl")
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
