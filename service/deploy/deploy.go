package deploy

import (
	"github.com/melbahja/goph"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/infra/typedef"
)

func Run(conf *typedef.DeployConf) error {
	defaultKey, err := getDefaultSSHPrivateKeyPath()
	if err != nil {
		return err
	}
	logrus.Debugln(defaultKey)
	auth, err := goph.Key(defaultKey, "")
	if err != nil {
		return err
	}

	client, err := goph.New(conf.User, conf.Server, auth)
	if err != nil {
		return err
	}
	// check service exist
	result, err := client.Run("systemctl status " + conf.Name)
	if err != nil {
		return err
	}
	logrus.Infoln(result)
	return nil
}
