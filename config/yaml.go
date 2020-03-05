package config

import (
	yaml2 "gopkg.in/yaml.v2"
	"hssh/common"
	"hssh/ssh"
	"io"
	"io/ioutil"
	"os"
)

// 读取配置文件
func ReadYamlConfig() (map[string]ssh.Config, bool) {
	if !common.CheckFileISExist(common.ConfigPath) {
		return make(map[string]ssh.Config), false
	}
	yaml, err := ioutil.ReadFile(common.ConfigPath)
	common.CheckErr(err)
	configs := make(map[string]ssh.Config)
	err = yaml2.Unmarshal(yaml, &configs)
	common.CheckErr(err)
	return configs, true
}

// 将sessions写入配置文件中
func WriteYamlConfig(sessions map[string]ssh.Config) {
	d, err := yaml2.Marshal(sessions)
	common.CheckErr(err)

	file, err := os.OpenFile(common.ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	common.CheckErr(err)
	defer file.Close()

	_, err = io.WriteString(file, string(d))
	common.CheckErr(err)
}

// 删除配置文件
func DelYamlFile() bool {
	if !common.CheckFileISExist(common.ConfigPath) {
		return true
	}
	err := os.Remove(common.ConfigPath)
	common.CheckErr(err)
	return true
}
