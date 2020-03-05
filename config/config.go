package config

import (
	"hssh/common"
	path2 "path"
)

var ProjectName string
var Version string
var YamlPath string

func init() {
	ProjectName = "hssh"
	Version = "0.0.11-beta"
	YamlPath = path2.Join(common.HomePath(), ".hssh.yaml")
}
