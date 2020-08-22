package payloads

import (
	"bytes"
	"text/template"
)

type HttpRequestBuilder struct {
	tmpl *template.Template
}

func (h HttpRequestBuilder) BuildPayload(values map[string]string) ([]byte, error) {
	var resolved bytes.Buffer

	if err := h.tmpl.Execute(&resolved, values); err != nil {
		return nil, err
	}

	return resolved.Bytes(), nil
}

func NewHttp(tmpl string) (PayloadBuilder, error) {
	parsedTmpl, err := template.New("HttpTemplate").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	return HttpRequestBuilder{parsedTmpl}, nil
}
