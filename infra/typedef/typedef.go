package typedef

type DeployConf struct {
	SystemdUnitConf
	SSHConf
}

type SystemdUnitConf struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Exec             string            `json:"exec,omitempty"`
	WorkingDirectory string            `json:"working_directory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
}

type SSHConf struct {
	Server string `json:"server,omitempty"`
	User   string `json:"user,omitempty"`
}

type LocalConf struct {
	Binary string            `json:"binary,omitempty"`
	Upload map[string]string `json:"upload,omitempty"`
}
