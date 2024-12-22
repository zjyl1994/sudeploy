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
	logrus.Infoln("Deploy", conf.Name, "to Remote", conf.Server, "User", conf.User, "Key", conf.Key, "KeyPass", conf.KeyPass != "")

	auth, err := goph.Key(conf.Key, conf.KeyPass)
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

	logrus.Infof("Remote exist %t,running %t,enabled %t\n", status.Exist, status.Running, status.Enabled)
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
	logrus.Debugf("Run script:\n%s\n", script)
	scriptName := "/tmp/sudeploy" + strconv.Itoa(rand.IntN(10000)) + ".sh"
	err := uploadTextFile(client, scriptName, script)
	if err != nil {
		return err
	}
	result, err := client.Run("bash " + scriptName)
	if err != nil {
		logrus.Errorf("Remote script:\n%s\n", string(result))
		return err
	}
	logrus.Infof("Remote script:\n%s\n", string(result))
	return nil
}

func uploadTextFile(c *goph.Client, filename, content string) error {
	sftp, err := c.NewSftp()
	if err != nil {
		return err
	}
	f, err := sftp.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	return err
}
