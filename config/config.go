package config

import (
	"hssh/common"
	path2 "path"
)

var ProjectName string
var Version string
var YamlPath string

type Server struct {
	User           string `yaml:"username"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"private_key"`
	KeyPassphrase  string `yaml:"key_passphrase"`
}

type Config struct {
	Servers map[string]Server `yaml:"servers,omitempty"`
}

func init() {
	ProjectName = "hssh"
	Version = "0.1.12-beta"
	YamlPath = path2.Join(common.HomePath(), ".hssh.yaml")
}
