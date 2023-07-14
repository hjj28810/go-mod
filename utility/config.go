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

// 从配置文件中载入json字符串
func LoadConfig[T any](path string, isDefault ...bool) T {
	var filePath string
	if len(isDefault) == 0 || isDefault[0] {
		filePath = GetCurrentPath()
	}
	buf, err := os.ReadFile(filePath + path)
	if err != nil {
		fmt.Println("load config conf failed: ", err)
	}
	return JsonBodyToObj[T](buf)
}

// 初始化 可以运行多次
func SetConfig[T any](path string, isDefault ...bool) T {
	return LoadConfig[T](path, isDefault...)
}

var instanceOnce sync.Once

// 初始化一次
func InitConfig[T any](path string, isDefault ...bool) T {
	var result T
	instanceOnce.Do(func() {
		result = LoadConfig[T](path, isDefault...)
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
