package cmd

import (
	"flag"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/service/deployconf"
	"github.com/zjyl1994/sudeploy/service/unitgen"
)

func Run() error {
	var confFile string
	flag.StringVar(&confFile, "conf", "sudeploy.json", "deploy config file")
	flag.Parse()

	absConfPath, err := filepath.Abs(confFile)
	if err != nil {
		return err
	}
	confFile = absConfPath
	conf, err := deployconf.Load(confFile)
	if err != nil {
		return err
	}

	logrus.Infoln("Deploying", conf.Name, "to", conf.Server)
	// TODO: check unit status in remote
	unitContent, err := unitgen.Gen(conf.SystemdUnitConf)
	if err != nil {
		return err
	}
	logrus.Infof("Systemd Unit:\n%s\n", unitContent)

	return nil
}
