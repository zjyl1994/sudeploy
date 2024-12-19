package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		logrus.Fatalln(err.Error())
	}
}
