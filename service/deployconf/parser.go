package deployconf

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/iancoleman/strcase"
	"github.com/zjyl1994/sudeploy/infra/typedef"
)

func Load(filename string) (*typedef.DeployConf, error) {
	bConf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result typedef.DeployConf
	err = json.Unmarshal(bConf, &result)
	if err != nil {
		return nil, err
	}
	result.Name = strcase.ToSnake(result.Name)
	if result.User == "" {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		result.User = u.Name
	}
	return &result, nil
}
