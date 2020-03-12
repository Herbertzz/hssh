package conf

import (
	"errors"
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
			return "", errors.New("private key does not exist")
		}
		return path, nil
	}
	fmt.Printf("The current path is not an absolute path, will use {home path}/%s\n\n", path)
	homePath, err := HomePath()
	if err != nil {
		return "", err
	}
	path = path2.Join(homePath, path)
	if !CheckFileISExist(path) {
		return "", errors.New("private key does not exist")
	}
	return path, nil
}

// ShowKeys 显示keys列表
func ShowKeys(configs Config) {
	index := 1
	for k, v := range configs.Keys {
		fmt.Printf("%02d. %s: %s\n", index, k, v)
		index++
	}
}
