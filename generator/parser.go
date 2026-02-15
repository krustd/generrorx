package generator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// EnumEntry 表示一个从 proto 解析出的错误码条目
type EnumEntry struct {
	Code     int    // 错误码数值
	Name     string // 枚举名，如 USER_NOT_FOUND
	Comment  string // 注释内容（内部错误的描述）
	HttpCode int    // @http=xxx，0 表示内部错误
	HttpMsg  string // 前端看到的消息
}

// IsInternal 判断是否为内部错误（没有 @http 标记）
func (e *EnumEntry) IsInternal() bool {
	return e.HttpCode == 0
}

// VarName 生成 Go 变量名，如 USER_NOT_FOUND → ErrUserNotFound
func (e *EnumEntry) VarName() string {
	return "Err" + toCamelCase(e.Name)
}

// ParseProto 解析 proto 文件中的 ErrorCode 枚举
func ParseProto(protoPath string) ([]*EnumEntry, error) {
	file, err := os.Open(protoPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 proto 文件: %w", err)
	}
	defer file.Close()

	var entries []*EnumEntry
	inEnum := false

	// 匹配 enum ErrorCode { 开始
	enumStartRe := regexp.MustCompile(`^\s*enum\s+ErrorCode\s*\{`)
	// 匹配 } 结束
	enumEndRe := regexp.MustCompile(`^\s*\}`)
	// 匹配枚举条目: NAME = 123; // 注释 @http=xxx
	entryRe := regexp.MustCompile(`^\s*(\w+)\s*=\s*(\d+)\s*;\s*(?://\s*(.*))?$`)
	// 匹配 @http=xxx
	httpRe := regexp.MustCompile(`(.+?)\s+@http=(\d+)\s*$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !inEnum {
			if enumStartRe.MatchString(line) {
				inEnum = true
			}
			continue
		}

		if enumEndRe.MatchString(line) {
			break
		}

		m := entryRe.FindStringSubmatch(line)
		if len(m) < 3 {
			continue
		}

		entry := &EnumEntry{
			Name: m[1],
			Code: toInt(m[2]),
		}

		comment := ""
		if len(m) == 4 {
			comment = strings.TrimSpace(m[3])
		}

		// 检查 @http 标记
		if hm := httpRe.FindStringSubmatch(comment); len(hm) == 3 {
			entry.HttpMsg = strings.TrimSpace(hm[1])
			entry.HttpCode = toInt(hm[2])
		} else {
			// 没有 @http → 内部错误，兜底 500 + "服务繁忙"
			entry.Comment = comment
			entry.HttpCode = 0
			entry.HttpMsg = ""
		}

		entries = append(entries, entry)
	}

	// 按 code 排序
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Code < entries[j].Code
	})

	return entries, nil
}

// ParsePbGo 解析已编译的 .pb.go 文件中的枚举映射（兼容旧流程）
func ParsePbGo(pbFile string) ([]*EnumEntry, error) {
	file, err := os.Open(pbFile)
	if err != nil {
		return nil, fmt.Errorf("无法打开 pb.go 文件: %w", err)
	}
	defer file.Close()

	enumMap := make(map[int]string)
	entryRe := regexp.MustCompile(`\s*(\d+):\s*"([^"]+)"`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := entryRe.FindStringSubmatch(line)
		if len(match) == 3 {
			code := toInt(match[1])
			name := match[2]
			enumMap[code] = name
		}
	}

	keys := make([]int, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var entries []*EnumEntry
	for _, code := range keys {
		entries = append(entries, &EnumEntry{
			Code: code,
			Name: enumMap[code],
		})
	}

	return entries, nil
}

// toCamelCase 将 SNAKE_CASE 转为 CamelCase
// USER_NOT_FOUND → UserNotFound
func toCamelCase(s string) string {
	parts := strings.Split(strings.ToLower(s), "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
