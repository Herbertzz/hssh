package conf

import (
	"fmt"
	"os"
	"os/user"
	path2 "path"
)

// CheckErr 检测是否有异常，如有则直接停止应用
func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// HomePath 返回当前用户的 home 路径
func HomePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}

// CheckFileISExist 判断文件是否存在
func CheckFileISExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// PrivateKeyPath 返回私钥的路径
func PrivateKeyPath(path string) (string, error) {
	if path2.IsAbs(path) {
		if !CheckFileISExist(path) {
			return "", fmt.Errorf("private key path(%s) not exist", path)
		}
		return path, nil
	}
	homePath, err := HomePath()
	if err != nil {
		return "", err
	}
	path = path2.Join(homePath, path)
	if !CheckFileISExist(path) {
		return "", fmt.Errorf("private key path(%s) not exist", path)
	}
	return path, nil
}

// ShowKeys 显示keys列表
func ShowKeys(configs Config) {
	keys := configs.SortPrivateKeys()
	var key Key
	for i, k := range keys {
		key = configs.Keys[k]
		fmt.Printf("%02d. %s: %s\n", i + 1, k, key.Path)
	}
}

// ShowServers 打印服务器列表
func ShowServers(config Config) {
	sessions := config.SortServerKeys()
	var authMethod string
	var session Server
	for i, k := range sessions {
		session = config.Servers[k]
		if session.AuthMethod == "password" {
			authMethod = "Password: " + session.Password
		} else if session.AuthMethod == "key" {
			authMethod = "Key: " + session.PrivateKey
		} else {
			authMethod = "undefined"
		}
		fmt.Printf("%02d. %s: %s@%s:%d(%s)\n", i + 1, k, session.User, session.Host, session.Port, authMethod)
	}
}