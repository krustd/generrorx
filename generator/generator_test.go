package generator

import (
	"os"
	"testing"
)

func TestParseProto(t *testing.T) {
	// 创建测试 proto 文件
	content := `syntax = "proto3";
package errorcode;
option go_package = "./errorcode";

enum ErrorCode {
    UNKNOWN             = 0;      // 未知错误

    USER_NOT_FOUND      = 10001;  // 用户不存在 @http=404
    INVALID_PARAM       = 10002;  // 参数错误 @http=400
    PASSWORD_WRONG      = 10003;  // 登录失败 @http=401
    ACCOUNT_LOCKED      = 10004;  // 登录失败 @http=401
    CAPTCHA_EXPIRED     = 10005;  // 登录失败 @http=401

    DB_CONNECT_FAILED   = 20001;  // 数据库连接失败
    DB_QUERY_TIMEOUT    = 20002;  // 数据库查询超时
    REDIS_TIMEOUT       = 20003;  // Redis超时
    RPC_CALL_FAILED     = 20004;  // RPC调用失败
}
`
	tmpFile := "/tmp/test_error.proto"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	entries, err := ParseProto(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(entries))
	}

	// 验证业务错误
	tests := []struct {
		idx      int
		name     string
		code     int
		httpCode int
		httpMsg  string
		internal bool
	}{
		{0, "UNKNOWN", 0, 0, "", true},
		{1, "USER_NOT_FOUND", 10001, 404, "用户不存在", false},
		{2, "INVALID_PARAM", 10002, 400, "参数错误", false},
		{3, "PASSWORD_WRONG", 10003, 401, "登录失败", false},
		{4, "ACCOUNT_LOCKED", 10004, 401, "登录失败", false},
		{5, "CAPTCHA_EXPIRED", 10005, 401, "登录失败", false},
		{6, "DB_CONNECT_FAILED", 20001, 0, "", true},
		{7, "DB_QUERY_TIMEOUT", 20002, 0, "", true},
		{8, "REDIS_TIMEOUT", 20003, 0, "", true},
		{9, "RPC_CALL_FAILED", 20004, 0, "", true},
	}

	for _, tt := range tests {
		e := entries[tt.idx]
		if e.Name != tt.name {
			t.Errorf("[%d] name: got %q, want %q", tt.idx, e.Name, tt.name)
		}
		if e.Code != tt.code {
			t.Errorf("[%d] code: got %d, want %d", tt.idx, e.Code, tt.code)
		}
		if e.HttpCode != tt.httpCode {
			t.Errorf("[%d] httpCode: got %d, want %d", tt.idx, e.HttpCode, tt.httpCode)
		}
		if e.HttpMsg != tt.httpMsg {
			t.Errorf("[%d] httpMsg: got %q, want %q", tt.idx, e.HttpMsg, tt.httpMsg)
		}
		if e.IsInternal() != tt.internal {
			t.Errorf("[%d] internal: got %v, want %v", tt.idx, e.IsInternal(), tt.internal)
		}
	}
}

func TestVarName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"USER_NOT_FOUND", "ErrUserNotFound"},
		{"UNKNOWN", "ErrUnknown"},
		{"DB_CONNECT_FAILED", "ErrDbConnectFailed"},
		{"RPC_CALL_FAILED", "ErrRpcCallFailed"},
		{"INVALID_PARAM", "ErrInvalidParam"},
	}

	for _, tt := range tests {
		e := &EnumEntry{Name: tt.name}
		if got := e.VarName(); got != tt.want {
			t.Errorf("VarName(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestGenerateDefault(t *testing.T) {
	entries := []*EnumEntry{
		{Code: 10001, Name: "USER_NOT_FOUND", HttpCode: 404, HttpMsg: "用户不存在"},
		{Code: 20001, Name: "DB_CONNECT_FAILED", Comment: "数据库连接失败"},
	}

	gen := &DefaultGenerator{}
	config := &Config{
		PackageName: "errorx",
		Framework:   "default",
	}

	files, err := gen.Generate(config, entries)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	// types.go 应该包含 Error 类型
	if files[0].Path != "./types.go" {
		t.Errorf("first file: got %q, want ./types.go", files[0].Path)
	}

	// errors_gen.go 应该包含变量定义
	errFile := files[1]
	if errFile.Path != "./errors_gen.go" {
		t.Errorf("second file: got %q, want ./errors_gen.go", errFile.Path)
	}

	content := errFile.Content
	if !contains(content, "ErrUserNotFound") {
		t.Error("missing ErrUserNotFound")
	}
	if !contains(content, "ErrDbConnectFailed") {
		t.Error("missing ErrDbConnectFailed")
	}
	if !contains(content, "404") {
		t.Error("missing http code 404")
	}
	if !contains(content, `"服务繁忙"`) {
		t.Error("missing default internal msg")
	}
}

func TestGenerateGoZero(t *testing.T) {
	entries := []*EnumEntry{
		{Code: 10001, Name: "USER_NOT_FOUND", HttpCode: 404, HttpMsg: "用户不存在"},
		{Code: 20001, Name: "DB_CONNECT_FAILED", Comment: "数据库连接失败"},
	}

	gen := &GoZeroGenerator{}
	config := &Config{
		PackageName: "errorx",
		Framework:   "gozero",
	}

	files, err := gen.Generate(config, entries)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}

	// 应该有 handler.go
	if files[2].Path != "./handler.go" {
		t.Errorf("third file: got %q, want ./handler.go", files[2].Path)
	}

	if !contains(files[2].Content, "ErrorHandler") {
		t.Error("handler.go missing ErrorHandler")
	}
	if !contains(files[2].Content, "logx.WithContext") {
		t.Error("handler.go missing logx usage")
	}
}

func TestGenerateGoFrame(t *testing.T) {
	entries := []*EnumEntry{
		{Code: 10001, Name: "USER_NOT_FOUND", HttpCode: 404, HttpMsg: "用户不存在"},
	}

	gen := &GoFrameGenerator{}
	config := &Config{
		PackageName: "errorx",
		Framework:   "goframe",
	}

	files, err := gen.Generate(config, entries)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}

	// types.go 应该有 ToGError
	if !contains(files[0].Content, "ToGError") {
		t.Error("types.go missing ToGError")
	}

	// middleware.go
	if files[2].Path != "./middleware.go" {
		t.Errorf("third file: got %q, want ./middleware.go", files[2].Path)
	}
	if !contains(files[2].Content, "ErrorMiddleware") {
		t.Error("middleware.go missing ErrorMiddleware")
	}
}

func contains(s, sub string) bool {
	return len(s) > 0 && len(sub) > 0 && len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
