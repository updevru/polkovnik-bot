package app

import (
	"bytes"
	"html/template"
)

type TemplateEngine struct {
	folder string
}

func NewTemplateEngine(folder string) *TemplateEngine {
	return &TemplateEngine{folder: folder}
}

func (t TemplateEngine) RenderString(name string, data interface{}) (string, error) {
	buf := &bytes.Buffer{}
	tpl, err := template.ParseFiles(t.folder + "/" + name)
	if err != nil {
		return "", err
	}

	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
