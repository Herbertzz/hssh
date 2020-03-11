package config

import (
	path2 "path"
)

var ProjectName = "hssh"
var Version = "0.4.14-beta"
var YamlPath string
var DefaultPrivateKey string

type Server struct {
	User          string `yaml:"username"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	AuthMethod    string `yaml:"auth_method"`
	Password      string `yaml:"password,omitempty"`
	PrivateKey    string `yaml:"private_key,omitempty"`
	KeyPassphrase string `yaml:"key_passphrase,omitempty"`
}

type Config struct {
	Keys    map[string]string `yaml:"keys"`
	Servers map[string]Server `yaml:"servers,omitempty"`
}

func init() {
	homePath, err := HomePath()
	CheckErr(err)

	YamlPath = path2.Join(homePath, ".hssh.yaml")
	DefaultPrivateKey = path2.Join(homePath, ".ssh/id_rsa")
}
