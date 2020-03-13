package conf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func getExecPath() string {
	// 获取当前执行程序的绝对路径
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

// in: error, expected: exit status 1
func TestCheckErr1(t *testing.T) {
	if os.Getenv("BE_CHECK_ERR") == "1" {
		CheckErr(errors.New("test"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErr1")
	cmd.Env = append(os.Environ(), "BE_CHECK_ERR=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// in: nil, expected: in
func TestCheckErr2(t *testing.T) {
	if os.Getenv("BE_CHECK_ERR") == "2" {
		CheckErr(nil)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErr2")
	cmd.Env = append(os.Environ(), "BE_CHECK_ERR=2")
	err := cmd.Run()
	if err == nil {
		return
	}

	t.Fatalf("process ran with err %v, want err nil", err)
}

func TestHomePath(t *testing.T) {
	actual, err := HomePath()
	if actual == "" {
		t.Errorf("HomePath() return empty string, expected: home path, Err: %v", err)
	}
}

func TestCheckFileISExist(t *testing.T) {
	var tests = []struct {
		in       string
		expected bool
	}{
		{getExecPath(), true},
		{"/tmp/not_exist_file_8dfk3d9", false},
	}

	for _, tt := range tests {
		actual := CheckFileISExist(tt.in)
		if actual != tt.expected {
			t.Errorf("CheckFileISExist(%s) = %t; expected %t", tt.in, actual, tt.expected)
		}
	}
}

func TestPrivateKeyPath(t *testing.T) {
	var tests = []struct {
		in       string
		expected string
	}{
		{getExecPath(), getExecPath()},
		{"/tmp/not_exist_file_8dfk3d9", ""},
		{DefaultPrivateKey, DefaultPrivateKey},
		{"no_exist_file_8mdf82li9", ""},
	}

	// 默认私钥不存在时，创建虚拟的私钥
	status := false
	if !CheckFileISExist(DefaultPrivateKey) {
		cmd := exec.Command("mkdir", "-p", DefaultPrivateKey)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			t.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		}
		status = true
	}
	// 完成时，删除该虚拟私钥
	defer func() {
		if status {
			cmd := exec.Command("rm", "-rf", DefaultPrivateKey)
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	for _, tt := range tests {
		actual, err := PrivateKeyPath(tt.in)
		if actual != tt.expected {
			t.Errorf("PrivateKeyPath(%s) = %s; expected %s; Error: %v", tt.in, actual, tt.expected, err)
		}
	}
}

func TestShowKeys(t *testing.T) {
	// 构造数据
	config := Config{
		Keys: map[string]Key{
			"default": {Path: "/Users/herbertzz/.ssh/id_rsa"},
			"xg":      {Path: "8msdfwr5544"},
		},
	}

	ShowKeys(config)
}
