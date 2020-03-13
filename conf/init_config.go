package conf

import (
	path2 "path"
)

// PROJECTNAME 项目名
const PROJECTNAME = "hssh"

// VERSION 版本号
const VERSION = "0.5.18-beta"

// YamlPath 配置文件路径
var YamlPath string

// DefaultPrivateKey 默认的私钥路径
var DefaultPrivateKey string

// Server 配置文件中的服务器配置结构
type Server struct {
	User          string `yaml:"username"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	AuthMethod    string `yaml:"auth_method"`
	Password      string `yaml:"password,omitempty"`
	PrivateKey    string `yaml:"private_key,omitempty"`
}

type Key struct {
	Path       string `yaml:"path"`
	Passphrase string `yaml:"passphrase,omitempty"`
}

// Config 配置文件结构
type Config struct {
	Keys    map[string]Key `yaml:"keys"`
	Servers map[string]Server `yaml:"servers,omitempty"`
}

func init() {
	homePath, err := HomePath()
	CheckErr(err)

	YamlPath = path2.Join(homePath, ".hssh.yaml")
	DefaultPrivateKey = path2.Join(homePath, ".ssh/id_rsa")
}
