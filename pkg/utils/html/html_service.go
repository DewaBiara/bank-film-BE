package html

import (
	"bytes"
)

type RenderService interface {
	GenerateHTMLDocument(docTemplate string, data *map[string]interface{}) (*bytes.Buffer, error)
}
