package cmd

import (
	"flag"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/service/deploy"
	"github.com/zjyl1994/sudeploy/service/deployconf"
)

func Run() error {
	var confFile string
	var debugMode bool
	flag.StringVar(&confFile, "conf", "sudeploy.json", "deploy config file")
	flag.BoolVar(&debugMode, "debug", false, "print more log")
	flag.Parse()

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
	
	absConfPath, err := filepath.Abs(confFile)
	if err != nil {
		return err
	}
	confFile = absConfPath
	conf, err := deployconf.Load(confFile)
	if err != nil {
		return err
	}

	return deploy.Run(conf)
}
