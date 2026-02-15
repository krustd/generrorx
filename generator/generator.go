package generator

import "fmt"

// OutputFile 表示一个待写入的输出文件
type OutputFile struct {
	Path    string
	Content string
}

// Config 生成器配置
type Config struct {
	PackageName    string // 生成代码所在的包名
	ImportPath     string // errorcode 包的 import 路径
	Framework      string // 目标框架: gozero, goframe, default
	DefaultHttpMsg string // 内部错误的默认前端消息
}

// Generator 代码生成器接口
type Generator interface {
	Generate(config *Config, entries []*EnumEntry) ([]*OutputFile, error)
}

// NewGenerator 根据 framework 名称创建对应的生成器
func NewGenerator(framework string) (Generator, error) {
	switch framework {
	case "default", "":
		return &DefaultGenerator{}, nil
	case "gozero":
		return &GoZeroGenerator{}, nil
	case "goframe":
		return &GoFrameGenerator{}, nil
	default:
		return nil, fmt.Errorf("不支持的框架: %s (可选: default, gozero, goframe)", framework)
	}
}
