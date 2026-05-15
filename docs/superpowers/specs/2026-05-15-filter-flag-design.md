# 设计文档：为所有 list 命令添加 `--filter` 标志

## 概述

为所有 `tapd <entity> list` 命令添加 `--filter "key=value"` 标志，支持 TAPD OpenAPI 的高级查询语法（LIKE、EQ、CONTAINS 等），让 AI Agent 能灵活组合过滤条件。

## 背景

TAPD OpenAPI 提供一套特殊查询语法，适用于所有字段：
- `field=LIKE<value>` — 模糊匹配
- `field=EQ<value>` — 精确匹配
- `field=NOT_EQ<value>` — 不等于
- `field=LIKE_OR<v1|v2>` — 模糊匹配多个值（OR）
- `field=CONTAINS<v1|v2|v3>` — 包含所有值（AND）
- `field=CONTAINS_OR<v1|v2|v3>` — 包含任一值（OR）
- `field=USER_OR<u1|u2>` — 多人查询（OR）
- 时间范围：`field=>2024-01-01`、`field=<2024-12-31`、`field=2024-01-01~2024-12-31`
- 多值：`field=v1|v2|v3`（OR）、`field=v1;v2`（AND）
- 否定：`field=<>"disabled"`

目前 CLI 的各 list 命令只暴露了少量硬编码的查询参数（如 `--status`、`--owner`），无法利用上述高级查询能力，也无法按自定义字段过滤。

## 设计方案

### 核心思路

CLI 层透传 `--filter` 参数到 SDK 的 `ToParams()` 返回的 map 中，不做任何校验或转换。利用 SDK 已有的透明 HTTP 参数传递机制（`map[string]string` → URL query string）。

### 实现细节

#### 1. 全局变量和解析函数

在 `internal/cmd/root.go` 中添加：

```go
var flagFilter []string

// parseFilters 将 ["key1=value1","key2=value2"] 解析为 map[string]string
func parseFilters(filters []string) map[string]string {
    if len(filters) == 0 {
        return nil
    }
    m := make(map[string]string, len(filters))
    for _, f := range filters {
        k, v, ok := strings.Cut(f, "=")
        if !ok || k == "" {
            continue
        }
        m[k] = v
    }
    return m
}
```

复用与 `parseCustomFields` 完全相同的解析逻辑（`key=value` 切割）。

#### 2. 每个命令注册标志

在每个 `xxxListCmd` 的 Flags 注册处添加：

```go
xxxListCmd.Flags().StringArrayVar(&flagFilter, "filter", nil, "高级过滤条件（可重复，格式：field=OP<value>）")
```

#### 3. 在 run 函数中合并参数

在每个 `runXxxList` 函数中，在调用 SDK 前将 filter 参数合并到请求 params：

```go
func runXxxList(cmd *cobra.Command, args []string) error {
    req := model.ListXxxRequest{...}
    params := req.ToParams()

    // 合并 filter 参数
    for k, v := range parseFilters(flagFilter) {
        params[k] = v
    }

    // 调用 SDK...
}
```

#### 4. 涉及的命令（14 个 list 命令）

| 命令 | 文件 | 预期收益 |
|------|------|----------|
| `story list` | `story.go` | 高 — 支持 custom_field_1-200 模糊/精确查询 |
| `bug list` | `bug.go` | 高 — 支持 custom_field_1-50 模糊/精确查询 |
| `task list` | `task.go` | 高 — 支持 custom_field_1-50 模糊/精确查询 |
| `iteration list` | `iteration.go` | 中 — 支持 custom_field_1-50 + 日期范围 |
| `release list` | `release.go` | 中 — 支持 name 模糊 + 日期范围 |
| `comment list` | `comment.go` | 低 — 时间范围/多ID 查询 |
| `timesheet list` | `timesheet.go` | 低 — 时间范围/多ID 查询 |
| `attachment list` | `attachment.go` | 低 — 文件名过滤 |
| `wiki list` | `wiki.go` | 中 — name 模糊匹配 |
| `tcase list` | `tcase.go` | 中 — 状态/优先级组合查询 |
| `category list` | `category.go` | 低 — name 过滤 |
| `workspace list` | `workspace.go` | 低 — 按状态/名称过滤 |
| `custom-field list` | `custom_field.go` | 低 — 按实体类型过滤 |
| `workitem-type list` | `custom_field.go` | 低 — 按状态过滤 |

#### 5. 用法示例

```bash
# 按自定义字段模糊搜索
tapd story list --filter "custom_field_one=LIKE<进度>"

# 按名称精确匹配
tapd story list --filter "name=EQ<登录功能>"

# 按时间范围查询
tapd bug list --filter "created=>2024-01-01" --filter "created=<2024-12-31"

# 组合多个条件
tapd task list --filter "status=CONTAINS_OR<开发中|测试中>" --filter "owner=USER_OR<张三|李四>"

# 与已有标志组合使用
tapd story list --owner zhangsan --filter "custom_field_two=EQ<高优先级>"
```

#### 6. 不做的事

- **不做参数校验**：CLI 不验证字段名、操作符或值格式，完全透传给 TAPD API，由 API 返回错误信息
- **不修改 SDK**：所有改动仅在 CLI 层
- **不替换现有标志**：`--status`、`--owner` 等现有标志保持不变，`--filter` 是增量的
- **不做帮助提示**：不做字段自动补全或可用字段提示

### 错误处理

如果 `--filter` 的值不包含 `=`，静默跳过（与 `--custom-field` 行为一致）。如果字段名或操作符不合法，TAPD API 会返回错误，CLI 直接透传该错误信息。

## 影响范围

- 修改文件：`root.go`（添加全局变量和解析函数）+ 所有 list 命令文件（注册标志 + 合并参数）
- 不影响现有功能
- 不需要修改 SDK
- 不需要修改测试框架（参数透传）
