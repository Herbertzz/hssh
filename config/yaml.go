package config

import (
	yaml2 "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	if !CheckFileISExist(YamlPath) {
		configs := InitYamlConfig()
		err := WriteYamlConfig(configs)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 初始化配置文件
func InitYamlConfig() Config {
	configs := Config{
		Keys: map[string]string{"default": DefaultPrivateKey},
	}
	return configs
}

// 读取配置文件
func ReadYamlConfig() (Config, error) {
	configs := Config{}
	yaml, err := ioutil.ReadFile(YamlPath)
	if err != nil {
		return configs, err
	}
	err = yaml2.Unmarshal(yaml, &configs)
	if err != nil {
		return configs, err
	}
	return configs, nil
}

// 写入配置文件
func WriteYamlConfig(configs Config) error {
	d, err := yaml2.Marshal(configs)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(YamlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func(f io.Closer) {
		if ferr := f.Close(); ferr != nil {
			log.Fatal(ferr)
		}
	}(file)

	_, err = io.WriteString(file, string(d))
	if err != nil {
		return err
	}
	return nil
}

// 删除配置文件
func DelYamlFile() error {
	if !CheckFileISExist(YamlPath) {
		return nil
	}
	err := os.Remove(YamlPath)
	if err != nil {
		return err
	}
	return nil
}
