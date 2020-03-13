package conf

import (
	"hssh/models"
	"os/exec"
	"reflect"
	"testing"
)

func TestInitYamlConfig(t *testing.T) {
	expected := Config{
		Keys:    map[string]string{"default": DefaultPrivateKey},
	}
	actual := InitYamlConfig()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("default conf error")
	}
}

func TestReadYamlConfig(t *testing.T) {
	_, err := ReadYamlConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteYamlConfig(t *testing.T) {
	configs, err := ReadYamlConfig()
	if err != nil {
		t.Fatal(err)
	}
	err = WriteYamlConfig(configs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelYamlFile(t *testing.T) {
	// 备份配置文件
	if models.CheckFileISExist(YamlPath) {
		cmd := exec.Command("cp", YamlPath, YamlPath + "_test_backup")
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
		// 完成时，备份拷回
		defer func() {
			cmd := exec.Command("mv", YamlPath + "_test_backup", YamlPath)
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	err := DelYamlFile()
	if err != nil {
		t.Fatal(err)
	}
}
