package conf

import (
	"reflect"
	"testing"
)

func TestInitYamlConfigWithNotExist(t *testing.T) {
	expected := Config{}
	if CheckFileISExist(DefaultPrivateKey) {
		err := moveTestFile(DefaultPrivateKey, DefaultPrivateKey+"_test_backup")
		if err != nil {
			t.Fatal(err)
		}
		// 完成时删除虚拟私钥
		defer func() {
			err := moveTestFile(DefaultPrivateKey+"_test_backup", DefaultPrivateKey)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	actual := InitYamlConfig()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("default conf error")
	}
}

func TestInitYamlConfigWithExist(t *testing.T) {
	expected := Config{
		Keys: map[string]Key{"default": {Path: DefaultPrivateKey}},
	}
	if !CheckFileISExist(DefaultPrivateKey) {
		err := createTestDir(DefaultPrivateKey)
		if err != nil {
			t.Fatal(err)
		}
		// 完成时删除虚拟私钥
		defer func() {
			err := destroyTestDir(DefaultPrivateKey)
			if err != nil {
				t.Fatal(err)
			}
		}()
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
	if CheckFileISExist(YamlPath) {
		err := copyTestFile(YamlPath, YamlPath+"_test_backup")
		if err != nil {
			t.Fatal(err)
		}
		// 完成时，备份拷回
		defer func() {
			err := moveTestFile(YamlPath+"_test_backup", YamlPath)
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
