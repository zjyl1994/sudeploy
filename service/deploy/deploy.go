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
	tmpBin := "/tmp/sudeploy" + strconv.Itoa(rand.IntN(10000))
	binaryPath := commandPathFromExec(conf.Exec)
	// upload new bin
	err := client.Upload(conf.Binary, tmpBin)
	if err != nil {
		return err
	}

	if state.Running {
		err := StopUnit(client, conf.Name)
		if err != nil {
			return err
		}
	}

	sftp, err := client.NewSftp()
	if err != nil {
		return err
	}

	err = sftp.Remove(binaryPath)
	if err != nil {
		return err
	}
	err = sftp.Rename(tmpBin, binaryPath)
	if err != nil {
		return err
	}
	return StartUnit(client, conf.Name)
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
	uploadTextFile(client, unitPath, unitFile)
	// reload daemon
	err = DaemonReload(client)
	if err != nil {
		return err
	}
	// start daemon
	err = StartUnit(client, conf.Name)
	if err != nil {
		return err
	}
	// enable daemon
	return EnableUnit(client, conf.Name)
}
