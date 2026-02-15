package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunBuild 编译 proto 文件
func RunBuild(protoPath string) {
	protoAbsPath, err := filepath.Abs(protoPath)
	if err != nil {
		fmt.Printf("❌ 无法解析 proto 路径: %s\n", err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat(protoAbsPath); os.IsNotExist(err) {
		fmt.Printf("❌ proto 文件不存在: %s\n", protoAbsPath)
		os.Exit(1)
	}

	protoDir := filepath.Dir(protoAbsPath)

	cmd := exec.Command("protoc",
		fmt.Sprintf("--go_out=%s", protoDir),
		fmt.Sprintf("--proto_path=%s", protoDir),
		protoAbsPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("🚀 编译 proto: %s\n", protoAbsPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ protoc 编译失败: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("✅ 已生成 .pb.go 文件")
}
