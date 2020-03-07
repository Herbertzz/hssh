package config

import (
	yaml2 "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

func init()  {
	if !CheckFileISExist(YamlPath) {
		InitYamlConfig()
	}
}

// 初始化配置文件
func InitYamlConfig() Config {
	keys := make(map[string]string)
	keys["default"] = DefaultPrivateKey
	configs := Config{
		Keys:    keys,
	}
	WriteProfile(configs)
	return configs
}

// 读取配置文件
func ReadYamlConfig() (Config, bool) {
	if !CheckFileISExist(YamlPath) {
		return InitYamlConfig(), false
	}
	yaml, err := ioutil.ReadFile(YamlPath)
	CheckErr(err)
	configs := Config{}
	err = yaml2.Unmarshal(yaml, &configs)
	CheckErr(err)
	return configs, true
}

// 将sessions写入配置文件中
func WriteYamlConfig(sessions map[string]Server) {
	configs, _ := ReadYamlConfig()
	configs.Servers = sessions

	WriteProfile(configs)
}

// 删除配置文件
func DelYamlFile() bool {
	if !CheckFileISExist(YamlPath) {
		return true
	}
	err := os.Remove(YamlPath)
	CheckErr(err)
	return true
}

// 写入配置文件
func WriteProfile(configs Config) {
	d, err := yaml2.Marshal(configs)
	CheckErr(err)

	file, err := os.OpenFile(YamlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	CheckErr(err)
	defer file.Close()

	_, err = io.WriteString(file, string(d))
	CheckErr(err)
}
