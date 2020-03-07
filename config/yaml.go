package config

import (
	yaml2 "gopkg.in/yaml.v2"
	"hssh/common"
	"io"
	"io/ioutil"
	"os"
)

// 读取配置文件
func ReadYamlConfig() (Config, bool) {
	if !common.CheckFileISExist(YamlPath) {
		return Config{}, false
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

	d, err := yaml2.Marshal(configs)
	common.CheckErr(err)

	file, err := os.OpenFile(YamlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	common.CheckErr(err)
	defer file.Close()

	_, err = io.WriteString(file, string(d))
	common.CheckErr(err)
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
