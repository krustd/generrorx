package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krustd/generrorx/generator"
)

// RunGenerate 生成错误码代码
func RunGenerate(opts GenOptions) {
	// 1. 选择解析方式
	var entries []*generator.EnumEntry
	var err error

	if opts.ProtoFile != "" {
		// 优先直接解析 proto 文件（推荐）
		entries, err = generator.ParseProto(opts.ProtoFile)
		if err != nil {
			fmt.Printf("❌ 解析 proto 文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📖 从 proto 文件解析到 %d 个错误码\n", len(entries))
	} else if opts.PbFile != "" {
		// 兼容旧流程：解析 .pb.go
		entries, err = generator.ParsePbGo(opts.PbFile)
		if err != nil {
			fmt.Printf("❌ 解析 pb.go 文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📖 从 pb.go 文件解析到 %d 个错误码\n", len(entries))
	} else {
		fmt.Println("❌ 请指定 --proto 或 --pbfile 参数")
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("⚠️  未找到任何错误码条目")
		return
	}

	// 2. 创建生成器
	gen, err := generator.NewGenerator(opts.Framework)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	// 3. 准备配置
	config := &generator.Config{
		PackageName:    opts.PackageName,
		ImportPath:     opts.ImportPath,
		Framework:      opts.Framework,
		DefaultHttpMsg: opts.DefaultHttpMsg,
	}

	// 4. 生成文件
	files, err := gen.Generate(config, entries)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		os.Exit(1)
	}

	// 5. 写入文件
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = "."
	}

	for _, f := range files {
		outPath := filepath.Join(outputDir, filepath.Base(f.Path))
		if err := os.WriteFile(outPath, []byte(f.Content), 0644); err != nil {
			fmt.Printf("❌ 写入文件失败 %s: %v\n", outPath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ 已生成: %s\n", outPath)
	}

	fmt.Printf("🎉 完成！框架: %s, 包名: %s, 共 %d 个错误码\n",
		opts.Framework, opts.PackageName, len(entries))
}

// GenOptions 生成命令的选项
type GenOptions struct {
	PackageName    string // 生成代码的包名
	ProtoFile      string // proto 文件路径（推荐）
	PbFile         string // .pb.go 文件路径（兼容旧流程）
	ImportPath     string // errorcode 包的 import 路径
	Framework      string // 目标框架
	OutputDir      string // 输出目录
	DefaultHttpMsg string // 内部错误的默认前端消息
}
