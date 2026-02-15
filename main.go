package main

import (
	"fmt"
	"os"

	"github.com/krustd/generrorx/cmd"

	"github.com/spf13/cobra"
)

var version = "v0.2.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:     "generrorx",
		Short:   "generrorx - 从 proto 生成多框架适配的错误码",
		Version: version,
		Long: `generrorx 是一个错误码生成工具，从 proto 文件定义错误码，
自动生成适配 go-zero / GoFrame / 通用 Go 项目的错误处理代码。

工作流程:
  1. generrorx create          # 创建 error.proto 模板
  2. 编辑 error.proto 添加错误码
  3. generrorx build           # 编译 proto（可选，兼容旧流程）
  4. generrorx gen -m <包名>   # 生成错误码代码

proto 标记规则:
  有 @http=xxx  → 业务错误，返回前端对应 HTTP 状态码 + 注释消息
  无 @http      → 内部错误，前端统一收到 500 + "服务繁忙"`,
	}

	// ==================== create ====================
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "创建 error.proto 模板文件",
		Run: func(c *cobra.Command, args []string) {
			cmd.RunCreate()
		},
	}

	// ==================== build ====================
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "编译 proto 文件（生成 .pb.go）",
		Run: func(c *cobra.Command, args []string) {
			path, _ := c.Flags().GetString("path")
			if path == "" {
				path = "./error.proto"
			}
			cmd.RunBuild(path)
		},
	}
	buildCmd.Flags().StringP("path", "p", "", "proto 文件路径（默认 ./error.proto）")

	// ==================== gen ====================
	var genCmd = &cobra.Command{
		Use:   "gen",
		Short: "生成错误码代码",
		Long: `从 proto 文件生成错误码代码，支持多种框架适配。

示例:
  # 通用模式（不依赖框架）
  generrorx gen -m myapp --proto ./error.proto

  # go-zero 适配（额外生成 handler.go）
  generrorx gen -m myapp --proto ./error.proto -f gozero

  # GoFrame 适配（额外生成 middleware.go）
  generrorx gen -m myapp --proto ./error.proto -f goframe

  # 指定输出目录
  generrorx gen -m myapp --proto ./error.proto -o ./errorx

  # 兼容旧流程（从 .pb.go 解析，不支持 @http 标记）
  generrorx gen -m myapp --pbfile ./errorcode/error.pb.go`,
		Run: func(c *cobra.Command, args []string) {
			opts := cmd.GenOptions{}
			opts.PackageName, _ = c.Flags().GetString("modelname")
			opts.ProtoFile, _ = c.Flags().GetString("proto")
			opts.PbFile, _ = c.Flags().GetString("pbfile")
			opts.ImportPath, _ = c.Flags().GetString("importpath")
			opts.Framework, _ = c.Flags().GetString("framework")
			opts.OutputDir, _ = c.Flags().GetString("output")
			opts.DefaultHttpMsg, _ = c.Flags().GetString("default-msg")

			if opts.Framework == "" {
				opts.Framework = "default"
			}
			if opts.ProtoFile == "" && opts.PbFile == "" {
				// 默认尝试 proto 文件
				if _, err := os.Stat("./error.proto"); err == nil {
					opts.ProtoFile = "./error.proto"
				} else {
					opts.PbFile = "./errorcode/error.pb.go"
				}
			}
			if opts.ImportPath == "" && opts.PackageName != "" {
				opts.ImportPath = opts.PackageName + "/errorcode"
			}

			cmd.RunGenerate(opts)
		},
	}
	genCmd.Flags().StringP("modelname", "m", "", "生成代码的包名（必填）")
	genCmd.MarkFlagRequired("modelname")
	genCmd.Flags().String("proto", "", "proto 文件路径（推荐，支持 @http 标记）")
	genCmd.Flags().String("pbfile", "", ".pb.go 文件路径（兼容旧流程）")
	genCmd.Flags().StringP("importpath", "i", "", "errorcode 包的 import 路径")
	genCmd.Flags().StringP("framework", "f", "default", "目标框架: default, gozero, goframe")
	genCmd.Flags().StringP("output", "o", "", "输出目录（默认当前目录）")
	genCmd.Flags().String("default-msg", "服务繁忙", "内部错误的默认前端消息")

	// ==================== auto ====================
	var autoCmd = &cobra.Command{
		Use:   "auto",
		Short: "一键完成：编译 proto + 生成错误码代码",
		Run: func(c *cobra.Command, args []string) {
			// Step 1: build proto
			path, _ := c.Flags().GetString("path")
			if path == "" {
				path = "./error.proto"
			}
			cmd.RunBuild(path)

			// Step 2: generate
			opts := cmd.GenOptions{}
			opts.PackageName, _ = c.Flags().GetString("modelname")
			opts.ProtoFile = path // 直接用 proto 文件解析
			opts.ImportPath, _ = c.Flags().GetString("importpath")
			opts.Framework, _ = c.Flags().GetString("framework")
			opts.OutputDir, _ = c.Flags().GetString("output")
			opts.DefaultHttpMsg, _ = c.Flags().GetString("default-msg")

			if opts.Framework == "" {
				opts.Framework = "default"
			}
			if opts.ImportPath == "" && opts.PackageName != "" {
				opts.ImportPath = opts.PackageName + "/errorcode"
			}

			cmd.RunGenerate(opts)
		},
	}
	autoCmd.Flags().StringP("path", "p", "", "proto 文件路径（默认 ./error.proto）")
	autoCmd.Flags().StringP("modelname", "m", "", "生成代码的包名（必填）")
	autoCmd.MarkFlagRequired("modelname")
	autoCmd.Flags().StringP("importpath", "i", "", "errorcode 包的 import 路径")
	autoCmd.Flags().StringP("framework", "f", "default", "目标框架: default, gozero, goframe")
	autoCmd.Flags().StringP("output", "o", "", "输出目录（默认当前目录）")
	autoCmd.Flags().String("default-msg", "服务繁忙", "内部错误的默认前端消息")

	rootCmd.AddCommand(createCmd, buildCmd, genCmd, autoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
