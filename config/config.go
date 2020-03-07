package config

import (
	"hssh/common"
	path2 "path"
)

var ProjectName string
var Version string
var YamlPath string
var DefaultPrivateKey string

type Server struct {
	User           string `yaml:"username"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	AuthMethod     string `yaml:"auth_method"`
	Password       string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty"`
	KeyPassphrase  string `yaml:"key_passphrase,omitempty"`
}

type Config struct {
	Keys    map[string]string `yaml:"keys"`
	Servers map[string]Server `yaml:"servers,omitempty"`
}

func init() {
	ProjectName = "hssh"
	Version = "0.2.13-beta"
	YamlPath = path2.Join(common.HomePath(), ".hssh.yaml")
	DefaultPrivateKey = path2.Join(common.HomePath(), ".ssh/id_rsa")
}
