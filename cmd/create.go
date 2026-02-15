package cmd

import (
	"fmt"
	"os"
)

// RunCreate 创建 error.proto 模板文件
func RunCreate() {
	const protoPath = "./error.proto"
	if _, err := os.Stat(protoPath); err == nil {
		fmt.Println("⚠️  error.proto 已存在，跳过")
		return
	}

	content := `syntax = "proto3";
package errorcode;
option go_package = "./errorcode";

// ErrorCode 错误码枚举
// 规则：
//   - 有 @http=xxx 标记的为业务错误，会返回给前端
//   - 没有 @http 标记的为内部错误，前端统一收到 500 + "服务繁忙"
//
// 示例：
//   USER_NOT_FOUND = 10001;  // 用户不存在 @http=404
//   DB_TIMEOUT     = 20001;  // 数据库超时
enum ErrorCode {
    UNKNOWN         = 0;      // 未知错误

    // ====== 业务错误（10000-19999）======
    // USER_NOT_FOUND  = 10001;  // 用户不存在 @http=404
    // INVALID_PARAM   = 10002;  // 参数错误 @http=400
    // UNAUTHORIZED    = 10003;  // 未授权 @http=401
    // FORBIDDEN       = 10004;  // 禁止访问 @http=403

    // ====== 内部错误（20000-29999）======
    // DB_CONNECT_FAILED  = 20001;  // 数据库连接失败
    // REDIS_TIMEOUT      = 20002;  // Redis超时
    // RPC_CALL_FAILED    = 20003;  // RPC调用失败
}
`
	err := os.WriteFile(protoPath, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ 已创建:", protoPath)
	fmt.Println("📝 请编辑 error.proto 添加你的错误码，然后运行:")
	fmt.Println("   generrorx build")
	fmt.Println("   generrorx gen -m <包名> -f <框架>")
}
