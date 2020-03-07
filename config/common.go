package config

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	path2 "path"
	"path/filepath"
)

// 检测是否有异常，如有则直接停止应用
func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// 返回当前用户的 home 路径
func HomePath() string {
	u, err := user.Current()
	CheckErr(err)
	return u.HomeDir
}

// 判断文件是否存在
func CheckFileISExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 删除当前执行程序
func DelCurrentApp() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	err := os.Remove(path)
	CheckErr(err)
}

// 返回私钥的路径
func PrivateKeyPath(path string) string {
	if path2.IsAbs(path) {
		return path
	}
	fmt.Printf("The current path is not an absolute path, will use {home path}/%s\n\n", path)
	return path2.Join(HomePath(), path)
}


// 显示keys列表
func ShowKeys(configs Config) {
	index := 1
	for k, v := range configs.Keys {
		fmt.Printf("%02d. %s: %s\n", index, k, v)
		index++
	}
}