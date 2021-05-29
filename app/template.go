package app

import (
	"bytes"
	"embed"
	"html/template"
)

type TemplateEngine struct {
	folder string
	files  embed.FS
}

func NewTemplateEngine(folder string, files embed.FS) *TemplateEngine {
	return &TemplateEngine{folder: folder, files: files}
}

func (t TemplateEngine) RenderString(name string, data interface{}) (string, error) {
	fileData, err := t.files.ReadFile(t.folder + "/" + name)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	tpl, err := template.New(name).Parse(string(fileData))
	if err != nil {
		return "", err
	}

	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
