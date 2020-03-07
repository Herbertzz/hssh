package config

import (
	yaml2 "gopkg.in/yaml.v2"
	"hssh/common"
	"io"
	"io/ioutil"
	"os"
)

func init()  {
	if !common.CheckFileISExist(YamlPath) {
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
	writeFile(configs)
	return configs
}

// 读取配置文件
func ReadYamlConfig() (Config, bool) {
	if !common.CheckFileISExist(YamlPath) {
		return InitYamlConfig(), false
	}
	yaml, err := ioutil.ReadFile(YamlPath)
	common.CheckErr(err)
	configs := Config{}
	err = yaml2.Unmarshal(yaml, &configs)
	common.CheckErr(err)
	return configs, true
}

// 将sessions写入配置文件中
func WriteYamlConfig(sessions map[string]Server) {
	configs, _ := ReadYamlConfig()
	configs.Servers = sessions

	writeFile(configs)
}

// 删除配置文件
func DelYamlFile() bool {
	if !common.CheckFileISExist(YamlPath) {
		return true
	}
	err := os.Remove(YamlPath)
	common.CheckErr(err)
	return true
}

// 写入配置文件
func writeFile(configs Config) {
	d, err := yaml2.Marshal(configs)
	common.CheckErr(err)

	file, err := os.OpenFile(YamlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	common.CheckErr(err)
	defer file.Close()

	_, err = io.WriteString(file, string(d))
	common.CheckErr(err)
}
