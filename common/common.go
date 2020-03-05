package common

import (
	"fmt"
	"os"
	"os/user"
	path2 "path"
)

var ConfigPath string

func init() {
	if ConfigPath == "" {
		ConfigPath = path2.Join(HomePath(), ".gssh.yaml")
	}
}

// 检测是否有异常，如有则直接停止应用
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
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

// 返回私钥的路径
func PrivateKeyPath(path string) string {
	if path2.IsAbs(path) {
		return path
	}
	return path2.Join(HomePath(), path)
}
