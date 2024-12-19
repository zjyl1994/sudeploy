package unitgen

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/zjyl1994/sudeploy/infra/typedef"
)

//go:embed systemd_unit.go.tmpl
var systemdUnitTemplate string

func Gen(conf typedef.SystemdUnitConf) (string, error) {
	tmpl, err := template.New("").Parse(systemdUnitTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, conf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
