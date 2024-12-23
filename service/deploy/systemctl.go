package deploy

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/melbahja/goph"
)

//go:embed script.go.tmpl
var deployScriptTmpl string

type UnitStatus struct {
	Exist   bool
	Running bool
	Enabled bool
}

func GetUnitStatus(client *goph.Client, unitName string) (UnitStatus, error) {
	var ret UnitStatus

	info, err := getUnitInfo(client, unitName)
	if err != nil {
		return ret, err
	}
	//logrus.Infoln(info)

	if state, ok := info["LoadState"]; ok {
		ret.Exist = strings.EqualFold(state, "loaded")
	}
	if state, ok := info["UnitFileState"]; ok {
		ret.Enabled = strings.EqualFold(state, "enabled")
	}
	if state, ok := info["ActiveState"]; ok {
		ret.Running = strings.EqualFold(state, "active")
	}

	return ret, nil
}

func getUnitInfo(client *goph.Client, unitName string) (map[string]string, error) {
	command := "systemctl show --no-pager " + unitName

	output, err := client.Run(command)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", command, err)
	}

	result := make(map[string]string)
	br := bufio.NewScanner(bytes.NewReader(output))
	br.Split(bufio.ScanLines)

	for br.Scan() {
		line := strings.TrimSpace(br.Text())
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func commandPathFromExec(exec string) string {
	parts := strings.Fields(exec)
	return parts[0]
}

type deployScriptParam struct {
	Name        string
	Install     bool
	Running     bool
	BinSrc      string
	BinDst      string
	WaitSeconds int
}

func GenDeployScript(param deployScriptParam) (string, error) {
	tmpl, err := template.New("").Parse(deployScriptTmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, param)
	return buf.String(), err
}
