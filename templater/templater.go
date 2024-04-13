package templater

import (
	"bytes"
	"text/template"
)

// RenderText uses data to render a text/template with the given name.
func RenderText(name, text string, data any, funcs template.FuncMap) ([]byte, error) {
	tpl, err := template.New(name).Funcs(funcs).Parse(text)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, data)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
