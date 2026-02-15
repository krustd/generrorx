# 更新日志

本文档记录了 generrorx 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 计划中
- 支持更多 Go Web 框架（Gin、Echo、Fiber 等）
- 添加 Web UI 配置界面
- 支持错误码自动递增
- 集成更多语言支持（Java、Python、TypeScript）
- 添加错误码版本管理
- 支持自定义错误码模板

## [0.2.0] - 2024-02-15

### 新增
- 🎉 初始版本发布
- ✨ 支持从 proto 文件生成错误码代码
- 🔄 多框架适配：go-zero、GoFrame、通用 Go 项目
- 🌐 内外错误分离：通过 `@http` 标记区分业务错误和系统错误
- 📝 标准 errors 兼容：支持 `errors.Is`、`errors.Unwrap`
- 🛠️ 完整的 CLI 工具：create、build、gen、auto 命令
- 📚 完善的文档和示例

### 支持的框架
- go-zero：生成 types.go、errors_gen.go、handler.go
- GoFrame：生成 types.go、errors_gen.go、middleware.go
- 通用模式：生成 types.go、errors_gen.go

### CLI 命令
- `generrorx create` - 创建 proto 模板
- `generrorx build` - 编译 proto 文件
- `generrorx gen` - 生成错误码代码
- `generrorx auto` - 一键完成编译和生成

---

## 版本说明

### 版本号格式

我们使用语义化版本号：`MAJOR.MINOR.PATCH`

- **MAJOR**：不兼容的 API 修改
- **MINOR**：向下兼容的功能性新增
- **PATCH**：向下兼容的问题修正

### 变更类型

- `新增` - 新功能
- `更改` - 对现有功能的变更
- `弃用` - 即将移除的功能
- `移除` - 已移除的功能
- `修复` - 问题修复
- `安全` - 安全相关的修复

---

## 贡献指南

如果您想为此项目做出贡献，请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

## 许可证

本项目采用 [GPL-3.0 许可证](LICENSE)。