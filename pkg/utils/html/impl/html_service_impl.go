package impl

import (
	"bytes"
	"html/template"

	"github.com/Budhiarta/bank-film-BE/pkg/utils/html"
)

type RenderServiceImpl struct {
}

func NewRenderServiceImpl() html.RenderService {
	return &RenderServiceImpl{}
}

func (*RenderServiceImpl) GenerateHTMLDocument(docTemplate string, data *map[string]interface{}) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles(docTemplate)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, *data); err != nil {
		return nil, err
	}

	return buf, nil
}
