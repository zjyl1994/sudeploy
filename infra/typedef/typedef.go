package typedef

type DeployConf struct {
	SystemdUnitConf
	SSHConf
	LocalConf
	WaitSeconds int `json:"wait_seconds,omitempty"`
}

type SystemdUnitConf struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Exec             string            `json:"exec,omitempty"`
	WorkingDirectory string            `json:"working_directory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
}

type SSHConf struct {
	Server  string `json:"server,omitempty"`
	Port    uint   `json:"port,omitempty"`
	User    string `json:"user,omitempty"`
	Key     string `json:"key,omitempty"`
	KeyPass string `json:"key_pass,omitempty"`
}

type LocalConf struct {
	Binary string            `json:"binary,omitempty"`
	Upload map[string]string `json:"upload,omitempty"`
}
