package template

import (
	"bytes"
	"text/template"

	"github.com/tonto/kit/goapp/pkg/bindata"
)

// Parse parses template
func Parse(tpl string, data interface{}) ([]byte, error) {
	d, err := bindata.Asset("data/template/" + tpl)
	if err != nil {
		return nil, err
	}

	t, err := template.New("tpl").Parse(string(d))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	err = t.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
