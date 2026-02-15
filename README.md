# generrorx

[![Go Version](https://img.shields.io/badge/Go-1.26+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-GPL%203.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/krustd/generrorx)](https://goreportcard.com/report/github.com/krustd/generrorx)
[![Release](https://img.shields.io/github/release/krustd/generrorx.svg)](https://github.com/krustd/generrorx/releases)
[![Build Status](https://img.shields.io/github/workflow/status/krustd/generrorx/CI)](https://github.com/krustd/generrorx/actions)

> 🚀 **从 proto 文件生成多框架适配的 Go 错误码代码** - 统一错误管理，提升开发效率

[English](#english) | [中文](#中文)

---

## 📖 目录

- [项目简介](#项目简介)
- [为什么选择 generrorx](#为什么选择-generrorx)
- [快速开始](#快速开始)
- [核心特性](#核心特性)
- [使用指南](#使用指南)
- [框架支持](#框架支持)
- [进阶用法](#进阶用法)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## 项目简介

**generrorx** 是一个强大的代码生成工具，专为 Go 项目设计，通过 Protocol Buffer 定义错误码，自动生成适配多种框架的错误处理代码。它解决了微服务架构中错误码管理混乱、前后端错误信息不一致的痛点。

### 🎯 解决的问题

- ❌ **错误码分散管理**：不同服务间错误码定义不一致
- ❌ **前后端信息不对称**：用户看到的错误信息与技术描述混淆
- ❌ **框架适配复杂**：不同框架需要不同的错误处理方式
- ❌ **维护成本高**：手动维护错误码容易出错且效率低下

### ✅ 我们的解决方案

- 📝 **单一数据源**：通过 proto 文件统一定义错误码
- 🔄 **多框架适配**：支持 go-zero、GoFrame、通用 Go 项目
- 🌐 **内外错误分离**：通过 `@http` 标记区分业务错误和系统错误
- 🎨 **标准化接口**：提供统一的错误处理 API

---

## 为什么选择 generrorx

| 特性 | generrorx | 传统方案 | 其他工具 |
|------|-----------|----------|----------|
| **Proto 定义** | ✅ 支持 | ❌ 不支持 | ⚠️ 部分支持 |
| **多框架适配** | ✅ 原生支持 | ❌ 需手动适配 | ⚠️ 有限支持 |
| **内外错误分离** | ✅ 标记区分 | ❌ 混合处理 | ❌ 不支持 |
| **标准 errors 兼容** | ✅ 完全兼容 | ⚠️ 部分兼容 | ⚠️ 部分兼容 |
| **一键生成** | ✅ 命令行工具 | ❌ 手动编写 | ⚠️ 配置复杂 |

---

## 快速开始

### 📦 安装

```bash
# 安装最新版本
go install github.com/krustd/generrorx@latest

# 或安装指定版本
go install github.com/krustd/generrorx@v0.2.0
```

### 🚀 三步上手

```bash
# 1. 创建 proto 模板
generrorx create

# 2. 编辑 error.proto 文件，定义你的错误码
vim error.proto

# 3. 生成代码（选择你的框架）
generrorx gen -m myapp -f gozero --proto ./error.proto
```

就这么简单！你的错误处理代码已经生成完毕 🎉

---

## 核心特性

### 🎯 Proto 定义错误码

使用 Protocol Buffer 作为单一数据源，前后端团队共享同一份错误定义：

```protobuf
syntax = "proto3";

package errorcode;

enum ErrorCode {
    // 默认错误
    UNKNOWN             = 0;      // 未知错误

    // 业务错误 - 有 @http 标记，返回前端
    USER_NOT_FOUND      = 10001;  // 用户不存在 @http=404
    INVALID_PARAM       = 10002;  // 参数错误 @http=400
    PASSWORD_WRONG      = 10003;  // 登录失败 @http=401
    ACCOUNT_LOCKED      = 10004;  // 账户被锁定 @http=401

    // 系统错误 - 无 @http 标记，前端统一收到 500
    DB_CONNECT_FAILED   = 20001;  // 数据库连接失败
    REDIS_TIMEOUT       = 20002;  // Redis 超时
    RPC_CALL_FAILED     = 20003;  // RPC 调用失败
}
```

### 🔄 多对一映射

多个业务错误码可以映射到同一个 HTTP 状态码：

```protobuf
USER_NOT_FOUND      = 10001;  // 用户不存在 @http=404
ORDER_NOT_FOUND     = 10005;  // 订单不存在 @http=404
PRODUCT_NOT_FOUND   = 10006;  // 商品不存在 @http=404
```

### 🌐 内外错误分离

- **有 `@http` 标记**：业务错误，返回前端对应 HTTP 状态码 + 注释消息
- **无 `@http` 标记**：系统错误，前端统一收到 500 + "服务繁忙"

---

## 使用指南

### 📋 CLI 参考

```bash
# 创建 proto 模板
generrorx create

# 编译 proto 文件（生成 .pb.go）
generrorx build [-p proto路径]

# 生成错误码代码
generrorx gen -m 包名 [选项]

# 一键完成：编译 proto + 生成错误码代码
generrorx auto -m 包名 [选项]
```

### ⚙️ 生成选项

| 选项 | 简写 | 描述 | 默认值 |
|------|------|------|--------|
| `--modelname` | `-m` | 生成代码的包名（必填） | - |
| `--framework` | `-f` | 目标框架: default, gozero, goframe | default |
| `--proto` | - | proto 文件路径（推荐） | ./error.proto |
| `--pbfile` | - | .pb.go 文件路径（兼容旧流程） | ./errorcode/error.pb.go |
| `--importpath` | `-i` | errorcode 包的 import 路径 | 包名/errorcode |
| `--output` | `-o` | 输出目录 | 当前目录 |
| `--default-msg` | - | 内部错误的默认前端消息 | "服务繁忙" |

### 📝 使用示例

```bash
# 通用模式（不依赖框架）
generrorx gen -m myapp --proto ./error.proto

# go-zero 适配（额外生成 handler.go）
generrorx gen -m myapp --proto ./error.proto -f gozero

# GoFrame 适配（额外生成 middleware.go）
generrorx gen -m myapp --proto ./error.proto -f goframe

# 指定输出目录
generrorx gen -m myapp --proto ./error.proto -o ./errorx

# 兼容旧流程（从 .pb.go 解析，不支持 @http 标记）
generrorx gen -m myapp --pbfile ./errorcode/error.pb.go
```

---

## 框架支持

### 🎯 go-zero 模式

生成文件：
```
types.go       # Error 类型定义
errors_gen.go  # 错误变量
handler.go     # ErrorHandler + OkHandler
```

使用方式：
```go
// main.go 注册错误处理
httpx.SetErrorHandlerCtx(errorx.ErrorHandler)

// service 层使用
func (s *UserService) Login(req *LoginReq) (*LoginResp, error) {
    user, err := s.repo.FindByPhone(req.Phone)
    if err != nil {
        return nil, errorx.ErrDbConnectFailed.Wrap(err)
        // 日志: [20001/DB_CONNECT_FAILED] dial tcp timeout
        // 前端: HTTP 500 {"code":500, "msg":"服务繁忙"}
    }
    if user == nil {
        return nil, errorx.ErrUserNotFound
        // 日志: [10001/USER_NOT_FOUND]
        // 前端: HTTP 404 {"code":404, "msg":"用户不存在"}
    }
    return &LoginResp{Token: genToken(user)}, nil
}
```

### 🎯 GoFrame 模式

生成文件：
```
types.go       # Error 类型定义 + ToGError
errors_gen.go  # 错误变量
middleware.go  # ErrorMiddleware
```

### 🎯 通用模式

生成文件：
```
types.go       # Error 类型定义
errors_gen.go  # 错误变量
```

---

## 进阶用法

### 🔧 错误包装与链式处理

```go
// 基础错误
err := errorx.ErrUserNotFound

// 包装错误（保留原始错误信息）
err = errorx.ErrDbConnectFailed.Wrap(originalErr)

// 链式包装
err = errorx.ErrRpcCallFailed.Wrap(originalErr).WithMsg("自定义消息")

// 错误判断
if errors.Is(err, errorx.ErrUserNotFound) {
    // 处理用户不存在的情况
}

// 获取错误码
if customErr, ok := err.(*errorx.Error); ok {
    fmt.Printf("错误码: %d\n", customErr.Code)
}
```

### 🌍 国际化支持

```go
// 在生成的代码基础上添加国际化支持
func (e *Error) Localize(lang string) string {
    switch lang {
    case "en":
        return e.MsgEn
    case "zh":
        return e.Msg
    default:
        return e.Msg
    }
}
```

### 📊 错误统计与监控

```go
// 添加错误统计功能
func (e *Error) WithMetrics() {
    metrics.Counter("error_count", map[string]string{
        "code": strconv.Itoa(e.Code),
        "name": e.Name,
    }).Inc()
}
```

---

## 贡献指南

我们欢迎所有形式的贡献！无论是提交 Issue、改进文档，还是提交代码。

### 🤝 如何贡献

1. **Fork** 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 **Pull Request**

### 📝 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/krustd/generrorx.git
cd generrorx

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 构建项目
go build -o generrorx main.go
```

### 🐛 报告 Bug

请使用 [Issues](https://github.com/krustd/generrorx/issues) 报告 Bug，并提供以下信息：

- 使用的 generrorx 版本
- 操作系统信息
- 重现步骤
- 预期行为与实际行为
- 相关的 proto 文件（如果适用）

### 💡 功能请求

欢迎提出新功能建议！请先查看 [Issues](https://github.com/krustd/generrorx/issues) 确认该功能尚未被提出。

---

## 社区与支持

- 📧 **邮箱**: krustd@github.com
- 💬 **讨论**: [GitHub Discussions](https://github.com/krustd/generrorx/discussions)
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/krustd/generrorx/issues)
- 📖 **文档**: [Wiki](https://github.com/krustd/generrorx/wiki)

---

## 许可证

本项目采用 [GPL-3.0 许可证](LICENSE)。

---

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=krustd/generrorx&type=Date)](https://star-history.com/#krustd/generrorx&Date)

---

## 🙏 致谢

感谢所有为 generrorx 做出贡献的开发者！

---

<div align="center">

**如果这个项目对你有帮助，请给我们一个 ⭐️**

Made with ❤️ by [krustd](https://github.com/krustd)

</div>

---

## English

### 🚀 generrorx - Generate Multi-Framework Adapted Go Error Code from Proto Files

A powerful code generation tool designed for Go projects that generates error handling code adapted to multiple frameworks through Protocol Buffer-defined error codes. It solves the pain points of error code management confusion and inconsistent frontend-backend error information in microservice architectures.

### 📦 Installation

```bash
go install github.com/krustd/generrorx@latest
```

### 🚀 Quick Start

```bash
# 1. Create proto template
generrorx create

# 2. Edit error.proto file
vim error.proto

# 3. Generate code
generrorx gen -m myapp -f gozero --proto ./error.proto
```

### 🎯 Key Features

- **Proto Definition**: Single source of truth for error codes
- **Multi-Framework Support**: go-zero, GoFrame, and generic Go projects
- **Internal/External Error Separation**: `@http` tags distinguish business from system errors
- **Standard errors Compatible**: Supports `errors.Is`, `errors.Unwrap`

### 📖 Documentation

For detailed documentation, please refer to the [Chinese section](#中文) above or visit our [Wiki](https://github.com/krustd/generrorx/wiki).

---

<div align="center">

**If this project helps you, please give us a ⭐️**

</div>