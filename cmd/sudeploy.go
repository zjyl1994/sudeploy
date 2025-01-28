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
	var skipVerify bool
	flag.StringVar(&confFile, "conf", "sudeploy.json", "deploy config file")
	flag.BoolVar(&debugMode, "debug", false, "print more log")
	flag.BoolVar(&skipVerify, "skip-verify", false, "skip ssh host key verify")
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

	if skipVerify && !conf.SkipVerify {
		conf.SkipVerify = skipVerify
	}

	return deploy.Run(conf)
}
