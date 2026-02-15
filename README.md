# generrorx

从 proto 文件生成多框架适配的 Go 错误码代码。

## 特性

- **Proto 定义错误码**：单一事实来源，前后端共享
- **`@http` 标记区分内外**：有标记的返回前端，无标记的仅日志
- **多对一映射**：多个后端错误码可映射到同一个 HTTP 状态码
- **多框架适配**：支持 go-zero、GoFrame、通用 Go 项目
- **标准 errors 兼容**：支持 `errors.Is`、`errors.Unwrap`

## 安装

```bash
go install github.com/krustd/generrorx@latest

```

## 快速开始

```bash
# 1. 创建 proto 模板
generrorx create

# 2. 编辑 error.proto
# 3. 生成代码（选择你的框架）
generrorx gen -m myapp -f gozero --proto ./error.proto
```

## Proto 标记规则

```protobuf
enum ErrorCode {
    UNKNOWN             = 0;      // 未知错误

    // 有 @http=xxx → 业务错误，返回前端
    USER_NOT_FOUND      = 10001;  // 用户不存在 @http=404
    INVALID_PARAM       = 10002;  // 参数错误 @http=400
    PASSWORD_WRONG      = 10003;  // 登录失败 @http=401
    ACCOUNT_LOCKED      = 10004;  // 登录失败 @http=401

    // 无 @http → 内部错误，前端统一收到 500 + "服务繁忙"
    DB_CONNECT_FAILED   = 20001;  // 数据库连接失败
    REDIS_TIMEOUT       = 20002;  // Redis超时
    RPC_CALL_FAILED     = 20003;  // RPC调用失败
}
```

注释 = 前端看到的消息（有 @http 时），后端日志描述（无 @http 时）。

## 生成结果

### 通用模式 (`-f default`)

```
types.go       # Error 类型定义
errors_gen.go  # 错误变量
```

### go-zero 模式 (`-f gozero`)

```
types.go       # Error 类型定义
errors_gen.go  # 错误变量
handler.go     # ErrorHandler + OkHandler
```

### GoFrame 模式 (`-f goframe`)

```
types.go       # Error 类型定义 + ToGError
errors_gen.go  # 错误变量
middleware.go  # ErrorMiddleware
```

## 使用示例

### go-zero 项目

**main.go 注册错误处理：**

```go
httpx.SetErrorHandlerCtx(errorx.ErrorHandler)
```

**service 层使用：**

```go
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
    if !checkPassword(user, req.Password) {
        return nil, errorx.ErrPasswordWrong
        // 日志: [10003/PASSWORD_WRONG]
        // 前端: HTTP 401 {"code":401, "msg":"登录失败"}
    }
    return &LoginResp{Token: genToken(user)}, nil
}
```

## CLI 参考

```
generrorx create                          # 创建 proto 模板
generrorx build [-p proto路径]             # 编译 proto
generrorx gen -m 包名 [选项]               # 生成错误码代码
generrorx auto -m 包名 [选项]              # 一键 build + gen

gen 选项:
  -m, --modelname    生成代码的包名（必填）
  -f, --framework    目标框架: default, gozero, goframe（默认 default）
      --proto        proto 文件路径（推荐）
      --pbfile       .pb.go 文件路径（兼容旧流程）
  -i, --importpath   errorcode 包的 import 路径
  -o, --output       输出目录（默认当前目录）
      --default-msg  内部错误的默认前端消息（默认"服务繁忙"）
```