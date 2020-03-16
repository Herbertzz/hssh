package conf

import (
	path2 "path"
	"sort"
)

// PROJECTNAME 项目名
const PROJECTNAME = "hssh"

// VERSION 版本号
const VERSION = "0.6.18-beta"

// YamlPath 配置文件路径
var YamlPath string

// DefaultPrivateKey 默认的私钥路径
var DefaultPrivateKey string

// Server 配置文件中的服务器配置结构
type Server struct {
	User       string `yaml:"username"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	AuthMethod string `yaml:"auth_method"`
	Password   string `yaml:"password,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty"`
}

// Key 配置文件中的私钥配置结构
type Key struct {
	Path       string `yaml:"path"`
	Passphrase string `yaml:"passphrase,omitempty"`
}

// Config 配置文件结构
type Config struct {
	Keys    map[string]Key    `yaml:"keys"`
	Servers map[string]Server `yaml:"servers,omitempty"`
}

// SortServerKeys 对服务器列表进行排序
func (config Config) SortServerKeys() []string {
	keys := make([]string, len(config.Servers))
	i := 0
	for k := range config.Servers {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortPrivateKeys 对私钥列表进行排序
func (config Config) SortPrivateKeys() []string {
	keys := make([]string, len(config.Keys))
	i := 0
	for k := range config.Keys {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func init() {
	homePath, err := HomePath()
	CheckErr(err)

	YamlPath = path2.Join(homePath, ".hssh.yaml")
	DefaultPrivateKey = path2.Join(homePath, ".ssh/id_rsa")
}
