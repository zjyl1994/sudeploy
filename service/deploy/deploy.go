package deploy

import (
	"io"
	"math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/melbahja/goph"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/sudeploy/infra/typedef"
	"github.com/zjyl1994/sudeploy/infra/vars"
	"github.com/zjyl1994/sudeploy/service/unitgen"
	"golang.org/x/crypto/ssh"
)

func Run(conf *typedef.DeployConf) error {
	port := uint(22)
	if conf.Port != 0 {
		port = conf.Port
	}

	logrus.Infoln("Deploy", conf.Name, "to Remote", conf.Server, conf.Port, "User", conf.User, "Key", conf.Key, "KeyPass", conf.KeyPass != "")

	auth, err := goph.Key(conf.Key, conf.KeyPass)
	if err != nil {
		return err
	}

	var hostVerifyFunc ssh.HostKeyCallback
	if conf.SkipVerify {
		hostVerifyFunc = func(string, net.Addr, ssh.PublicKey) error {
			return nil
		}
	} else {
		callback, err := goph.DefaultKnownHosts()
		if err != nil {
			return err
		}
		hostVerifyFunc = callback
	}

	client, err := goph.NewConn(&goph.Config{
		User:     conf.User,
		Addr:     conf.Server,
		Port:     port,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: hostVerifyFunc,
	})
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
	err := uploadBinaryFile(client, conf.Binary, tmpBin)
	if err != nil {
		return err
	}
	// run deploy script
	script, err := GenDeployScript(deployScriptParam{
		Name:        conf.Name,
		Running:     state.Running,
		BinSrc:      tmpBin,
		BinDst:      binaryPath,
		WaitSeconds: conf.WaitSeconds,
	})
	if err != nil {
		return err
	}
	return runDeployScript(client, script)
}

func installBinary(client *goph.Client, conf *typedef.DeployConf) error {
	binaryPath := commandPathFromExec(conf.Exec)
	for localPath, remotePath := range conf.Upload {
		if err := uploadBinaryFile(client, localPath, remotePath); err != nil {
			return err
		}
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
	err = uploadBinaryFile(client, conf.Binary, tmpBin)
	if err != nil {
		return err
	}
	// run deploy script
	script, err := GenDeployScript(deployScriptParam{
		Name:        conf.Name,
		Install:     true,
		BinSrc:      tmpBin,
		BinDst:      binaryPath,
		WaitSeconds: conf.WaitSeconds,
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

func uploadBinaryFile(c *goph.Client, localPath, remotePath string) (err error) {
	local, err := os.Open(localPath)
	if err != nil {
		return
	}
	defer local.Close()

	fi, err := local.Stat()
	if err != nil {
		return
	}

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Create(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	bar := progressbar.DefaultBytes(
		fi.Size(),
		filepath.Base(localPath),
	)

	_, err = io.Copy(io.MultiWriter(remote, bar), local)
	return
}
