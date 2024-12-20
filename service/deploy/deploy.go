package deploy

import (
	"math/rand/v2"
	"path/filepath"
	"strconv"

	"github.com/melbahja/goph"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/infra/typedef"
	"github.com/zjyl1994/sudeploy/infra/vars"
	"github.com/zjyl1994/sudeploy/service/unitgen"
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
	// check service status
	status, err := GetUnitStatus(client, conf.Name)
	if err != nil {
		return err
	}

	if status.Exist {
		return updateBinary(client, conf, status)
	} else {
		return installBinary(client, conf)
	}
}

func updateBinary(client *goph.Client, conf *typedef.DeployConf, state UnitStatus) error {
	tmpBin := "/tmp/sudeploy" + strconv.Itoa(rand.IntN(10000)) + ".bin"
	binaryPath := commandPathFromExec(conf.Exec)
	// upload new bin
	err := client.Upload(conf.Binary, tmpBin)
	if err != nil {
		return err
	}
	// run deploy script
	script, err := GenDeployScript(deployScriptParam{
		Name:    conf.Name,
		Running: state.Running,
		BinSrc:  tmpBin,
		BinDst:  binaryPath,
	})
	if err != nil {
		return err
	}
	return runDeployScript(client, script)
}

func installBinary(client *goph.Client, conf *typedef.DeployConf) error {
	binaryPath := commandPathFromExec(conf.Exec)
	err := client.Upload(conf.Binary, binaryPath)
	if err != nil {
		return err
	}
	// upload unit file
	unitFile, err := unitgen.Gen(conf.SystemdUnitConf)
	if err != nil {
		return err
	}
	unitPath := filepath.Join(vars.DefaultSystemdUnitPath, conf.Name+".service")
	err = uploadTextFile(client, unitPath, unitFile)
	if err != nil {
		return err
	}
	// upload new bin
	tmpBin := "/tmp/sudeploy" + strconv.Itoa(rand.IntN(10000)) + ".bin"
	err = client.Upload(conf.Binary, tmpBin)
	if err != nil {
		return err
	}
	// run deploy script
	script, err := GenDeployScript(deployScriptParam{
		Name:    conf.Name,
		Install: true,
		BinSrc:  tmpBin,
		BinDst:  binaryPath,
	})
	if err != nil {
		return err
	}
	return runDeployScript(client, script)
}

func runDeployScript(client *goph.Client, script string) error {
	scriptName := "/tmp/sudeploy" + strconv.Itoa(rand.IntN(10000)) + ".sh"
	err := uploadTextFile(client, scriptName, script)
	if err != nil {
		return err
	}
	result, err := client.Run("bash " + scriptName)
	if err != nil {
		logrus.Errorln("Remote script:", result)
		return err
	}
	logrus.Infoln("Remote script:", result)
	return nil
}
