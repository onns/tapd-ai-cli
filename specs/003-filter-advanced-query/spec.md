# 功能规范: 所有 list 命令增加 --filter 高级过滤

**功能分支**: `003-filter-advanced-query`
**创建时间**: 2026-05-15
**状态**: 已实现
**输入**: 用户描述: "增加自定义字段以及已有字段的高级别的过滤能力，基于 openapi 的高级查找能力"
**设计文档**: `docs/superpowers/specs/2026-05-15-filter-flag-design.md`

## 用户场景与测试 *(必填)*

### 用户故事 1 - AI Agent 按自定义字段模糊搜索需求 (优先级: P1)

作为 AI Agent，我需要按自定义字段（如 custom_field_one）模糊搜索需求列表，以便快速定位特定类型的需求。

**验收场景**:
1. **给定** 工作区存在自定义字段 `custom_field_one`，**当** 执行 `tapd story list --filter "custom_field_one=LIKE<进度>"` 时，**那么** 返回匹配的需求数组
2. **给定** 多个过滤条件，**当** 执行 `tapd story list --filter "name=LIKE<登录>" --filter "status=EQ<已实现>"` 时，**那么** 返回同时满足两个条件的需求数组

---

### 用户故事 2 - AI Agent 按时间范围组合查询 (优先级: P1)

作为 AI Agent，我需要按时间范围组合查询缺陷列表，以便统计特定时间段内创建的缺陷。

**验收场景**:
1. **给定** 工作区存在缺陷数据，**当** 执行 `tapd bug list --filter "created=>2024-01-01" --filter "created=<2024-12-31"` 时，**那么** 返回该时间段内创建的缺陷数组

---

### 用户故事 3 - 所有 list 命令统一支持 --filter (优先级: P2)

作为 AI Agent，我希望所有 list 命令都支持 `--filter` 标志，以便在任意实体类型的列表查询中使用高级过滤语法。

**验收场景**:
1. **给定** 任意 list 命令（story/bug/task/iteration/release/comment/timesheet/attachment/wiki/tcase/category/custom-field/workitem-type），**当** 执行带 `--filter` 参数的命令时，**那么** 参数被透传到 TAPD API，返回过滤后的结果

## 技术方案

### 核心设计

- **CLI 层透传**：`--filter "key=value"` 解析后合并到 SDK 的 `ToParams()` 返回 map，不做校验
- **SDK 最小改动**：暴露 `DoGet` 公开方法（绕过 struct 到 map 的内部转换）和 `ParseList` 泛型函数（处理 TAPD 包装响应）
- **泛型 helper**：`listWithFilters[T any]` 统一处理参数合并、API 调用和响应解析

### TAPD OpenAPI 特殊查询语法

| 操作符 | 含义 | 示例 |
|--------|------|------|
| `LIKE<value>` | 模糊匹配 | `name=LIKE<登录>` |
| `EQ<value>` | 精确匹配 | `status=EQ<已实现>` |
| `NOT_EQ<value>` | 不等于 | `status=NOT_EQ<已关闭>` |
| `LIKE_OR<v1\|v2>` | 模糊匹配多值（OR） | `name=LIKE_OR<登录\|注册>` |
| `CONTAINS<v1\|v2>` | 包含所有值（AND） | `label=CONTAINS<前端\|高优>` |
| `CONTAINS_OR<v1\|v2\|v3>` | 包含任一值（OR） | `status=CONTAINS_OR<开发中\|测试中>` |
| `USER_OR<u1\|u2>` | 多人查询（OR） | `owner=USER_OR<张三\|李四>` |
| `>` / `<` | 大于/小于（时间/数值） | `created=>2024-01-01` |
| `~` | 时间范围 | `created=2024-01-01~2024-12-31` |

### 涉及改动的文件

**SDK 改动**（tapd-sdk-go 仓库）：
- `client.go`：新增公开 `DoGet` 方法
- `parse.go`：新增公开 `ParseList` 泛型函数

**CLI 改动**（tapd-ai-cli 仓库）：
- `go.mod`：添加 `replace` 指令指向本地 SDK
- `internal/cmd/root.go`：新增 `flagFilter`、`listWithFilters` helper
- `internal/cmd/skill_template.md`：新增高级过滤章节
- 12 个 list 命令文件：注册 `--filter` 标志，改用 `listWithFilters`
- `README.md`、`CODEBUDDY.md`：文档更新

### 契约变更（相对于 001-mvp-tapd-cli）

详见 [contracts/filter-contract.md](contracts/filter-contract.md)
