package conf

import (
	"bytes"
	"html/template"
)

const ext_tmpl = `
[{{ .ContextName }}]\n
exten = s,1,Answer()
{{range $sound := .Sounds}}
same = n,Wait(1)
same = n,Playback({{ $sound }})
{{end}}
`

var confs []ExtConfig

type ExtConfig struct {
	ContextName string
	Sounds      []string
}

func (c *ExtConfig) ToString() string {
	var data bytes.Buffer
	t := template.Must(template.New("ext").Parse(ext_tmpl))

	t.Execute(&data, c)

	return data.String()
}
