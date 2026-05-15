---
name: tapd
description: TAPD操作，涉及需求，缺陷，任务，Wiki等。tapd.cn快捷查询：tapd url <url>
---

面向 AI Agent 的 TAPD 命令行工具。通过 `tapd` 命令与 TAPD 平台交互，所有输出针对最小 token 消耗优化。

## 安装

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

## 认证

```bash
# Access Token（推荐）
export TAPD_ACCESS_TOKEN=<your_token>

# 或交互式登录持久化凭据
tapd auth login
```

凭据优先级：CLI flags > 环境变量 > `./.tapd.json` > `~/.tapd.json`

## 高级过滤（--filter）

所有 `list` 命令均支持 `--filter` 标志，可重复使用，直接透传 TAPD OpenAPI 的高级查询语法：

```bash
tapd story list --filter "name=LIKE<登录>" --filter "status=EQ<已实现>"
tapd bug list --filter "created=>2024-01-01" --filter "custom_field_one=EQ<高优先级>"
tapd task list --filter "owner=USER_OR<张三|李四>"
```

### 支持的操作符

| 操作符 | 含义 | 示例 |
|--------|------|------|
| `LIKE<value>` | 模糊匹配 | `name=LIKE<登录>` |
| `EQ<value>` | 精确匹配 | `status=EQ<已实现>` |
| `NOT_EQ<value>` | 不等于 | `status=NOT_EQ<已关闭>` |
| `LIKE_OR<v1\|v2>` | 模糊匹配多个值（OR） | `name=LIKE_OR<登录\|注册>` |
| `CONTAINS<v1\|v2>` | 包含所有值（AND） | `label=CONTAINS<前端\|高优>` |
| `CONTAINS_OR<v1\|v2\|v3>` | 包含任一值（OR） | `status=CONTAINS_OR<开发中\|测试中>` |
| `USER_OR<u1\|u2>` | 多人查询（OR） | `owner=USER_OR<张三\|李四>` |
| `>` / `<` | 大于/小于（时间/数值） | `created=>2024-01-01` |
| `~` | 时间范围 | `created=2024-01-01~2024-12-31` |
| `\|` | 多值 OR | `status=开发中\|测试中` |
| `<>` | 不等于（简写） | `status=<>已关闭` |

### 适用字段

- 标准字段（name/title/status/owner/created 等）
- 自定义字段（custom_field_one、custom_field_two 等，用 `custom-field list --entity-type stories` 查看可用字段）

`--filter` 与已有标志（`--status`、`--owner` 等）可组合使用，参数会合并传递给 API。

## 命令参考

{{COMMAND_REFERENCE}}
