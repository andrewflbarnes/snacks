package loris

import (
	"bytes"
	"text/template"
)

var (
	httpTmpl = `POST {{.Endpoint}} HTTP/1.1
Content-Length: {{.Length}}
Content-Type: {{.ContentType}}
User-Agent: snacks

{{.Body}}


`
)

type Loris struct {
	Tmpl *template.Template
}

type LorisVals struct {
	Endpoint    string
	Length      int
	ContentType string
	Body        string
}

func (l Loris) Build(vals LorisVals) string {
	var resolved bytes.Buffer
	if err := l.Tmpl.Execute(&resolved, vals); err != nil {
		panic(err)
	}

	return resolved.String()
}

func New() Loris {
	tmpl, err := template.New("default").Parse(httpTmpl)

	if err != nil {
		panic(err)
	}

	return Loris{
		Tmpl: tmpl,
	}
}
