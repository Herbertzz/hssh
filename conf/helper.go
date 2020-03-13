package conf

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// 测试用例用，用于创建测试用临时目录
func createTestDir(path string) error {
	cmd := exec.Command("mkdir", "-p", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}

// 测试用例用，用于删除测试用临时目录
func destroyTestDir(path string) error {
	cmd := exec.Command("rm", "-rf", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}
	return nil
}

// 测试用例用， 用于移动文件
func moveTestFile(path string, dst string) error {
	cmd := exec.Command("mv", path, dst)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}

// 测试用例用， 用于移动文件
func copyTestFile(path string, dst string) error {
	cmd := exec.Command("cp", path, dst)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}
