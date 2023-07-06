package utility

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type configType int

const (
	DevConfig configType = iota
	TestConfig
	ProdConfig
	OtherConfig
)

func (c configType) ToPath(path string) string {
	switch c {
	case DevConfig:
		return path + "/config_dev.json"
	case TestConfig:
		return path + "/config_test.json"
	case ProdConfig:
		return path + "/config_prod.json"
	default:
		return path
	}
}

var instanceOnce sync.Once

// 从配置文件中载入json字符串
func LoadConfig[T any](path string) T {
	buf, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("load config conf failed: ", err)
	}
	return JsonBodyToObj[T](buf)
}

// 初始化 可以运行多次
func SetConfig[T any](path string) T {
	return LoadConfig[T](path)
}

// 初始化，只运行一次
func InitConfig[T any](path string) T {
	filePath := GetCurrentPath()
	var result T
	instanceOnce.Do(func() {
		result = LoadConfig[T](filePath + path)
	})
	return result
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetCurrentPath() string {
	binary, err := os.Executable()
	rootPath := filepath.Dir(filepath.Dir(binary))
	CheckErr(err)
	return rootPath
}
