# 贡献指南

感谢您对 generrorx 项目的关注！我们欢迎所有形式的贡献，无论是报告问题、改进文档，还是提交代码。

## 🤝 贡献方式

### 报告 Bug

如果您发现了 Bug，请通过 [Issues](https://github.com/krustd/generrorx/issues) 报告，并提供以下信息：

- **Bug 描述**：清晰描述遇到的问题
- **重现步骤**：详细的重现步骤
- **预期行为**：您期望发生什么
- **实际行为**：实际发生了什么
- **环境信息**：
  - 操作系统
  - Go 版本
  - generrorx 版本
- **相关文件**：如果有相关的 proto 文件或生成的代码，请提供

### 功能请求

如果您有新功能的想法，请：

1. 先查看 [Issues](https://github.com/krustd/generrorx/issues) 确认该功能尚未被提出
2. 创建新的 Issue，详细描述：
   - 功能描述
   - 使用场景
   - 预期的 API 设计
   - 可能的实现方案

### 提交代码

#### 开发环境设置

1. **Fork 仓库**
   ```bash
   # 在 GitHub 上 Fork 仓库，然后克隆到本地
   git clone https://github.com/YOUR_USERNAME/generrorx.git
   cd generrorx
   ```

2. **添加上游仓库**
   ```bash
   git remote add upstream https://github.com/krustd/generrorx.git
   ```

3. **安装依赖**
   ```bash
   go mod tidy
   ```

4. **运行测试**
   ```bash
   go test ./...
   ```

#### 开发流程

1. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

2. **进行开发**
   - 遵循现有的代码风格
   - 添加必要的测试
   - 更新相关文档

3. **运行测试和检查**
   ```bash
   # 运行所有测试
   go test ./...
   
   # 运行代码格式化
   go fmt ./...
   
   # 运行 go vet
   go vet ./...
   
   # 构建
   go build -o generrorx main.go
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 填写 PR 模板
   - 等待代码审查

#### 代码规范

- **命名规范**：
  - 使用驼峰命名法
  - 包名使用小写字母
  - 常量使用大写字母加下划线

- **注释规范**：
  - 公开的函数必须有注释
  - 复杂逻辑需要添加行内注释
  - 使用 godoc 格式

- **错误处理**：
  - 不要忽略错误
  - 使用项目定义的错误类型
  - 提供有意义的错误信息

#### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

类型包括：
- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式化（不影响功能）
- `refactor`: 重构代码
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(generator): add support for gin framework

- Add gin framework adapter
- Update documentation
- Add integration tests

Closes #123
```

## 📝 文档贡献

文档是项目的重要组成部分，我们欢迎以下类型的文档贡献：

- **README 改进**：让项目更容易理解和使用
- **API 文档**：完善函数和类型的文档
- **教程**：编写使用教程和最佳实践
- **示例**：提供更多使用示例
- **翻译**：将文档翻译成其他语言

## 🏷️ 标签和分类

在提交 Issue 或 PR 时，我们会使用以下标签：

### Issue 标签

- `bug`: Bug 报告
- `enhancement`: 功能增强
- `documentation`: 文档相关
- `good first issue`: 适合新贡献者的问题
- `help wanted`: 需要帮助的问题

### PR 标签

- `ready for review`: 准备审查
- `work in progress`: 进行中的工作
- `do not merge`: 请勿合并
- `needs tests`: 需要添加测试

## 🎯 贡献重点

我们特别欢迎以下方面的贡献：

### 高优先级

1. **新框架支持**：添加对更多 Go Web 框架的支持
2. **错误码管理**：改进错误码管理和版本控制
3. **性能优化**：提升代码生成和运行时性能
4. **测试覆盖率**：提高测试覆盖率，特别是边界情况

### 中优先级

1. **文档完善**：改进现有文档，添加更多示例
2. **工具集成**：与 CI/CD 工具的集成
3. **国际化**：支持多语言错误消息
4. **插件系统**：支持自定义插件扩展功能

### 低优先级

1. **UI 界面**：开发 Web UI 配置界面
2. **多语言支持**：支持生成其他语言的代码
3. **模板系统**：允许用户自定义代码模板

## 🚀 发布流程

我们使用语义化版本控制 (SemVer)：

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

发布流程：

1. 更新版本号（`main.go` 中的 `version` 变量）
2. 更新 `CHANGELOG.md`
3. 创建 Git 标签
4. 创建 GitHub Release
5. 发布到 Go Modules

## 📞 联系方式

如果您有任何问题或建议，可以通过以下方式联系我们：

- **GitHub Issues**: [提交问题](https://github.com/krustd/generrorx/issues)
- **GitHub Discussions**: [参与讨论](https://github.com/krustd/generrorx/discussions)
- **Email**: krustd@github.com

## 🙏 致谢

感谢所有为 generrorx 做出贡献的开发者！您的贡献让这个项目变得更好。

---

再次感谢您的贡献！每一个贡献都是宝贵的。

---

## 许可证

本项目采用 [GPL-3.0 许可证](LICENSE)。在贡献代码时，您同意您的贡献将在同一许可证下发布。