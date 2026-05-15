# 契约变更: --filter 标志

**基线**: `specs/001-mvp-tapd-cli/contracts/cli-commands.md`

本文档描述 003 功能对 001 契约的增量变更。所有变更仅涉及 `list` 子命令新增 `--filter` 标志。

## 新增全局 flag

### --filter

| 属性 | 值 |
|------|-----|
| 类型 | `stringArray`（可重复） |
| 适用范围 | 所有 `list` 子命令 |
| 格式 | `field=OP<value>` |
| 校验 | 无，完全透传给 TAPD API |

## 受影响的命令

以下命令新增 `--filter` 标志，与已有标志可组合使用：

### tapd story list

```diff
-tapd story list [--status <status>] [--owner <owner>] [--iteration-id <id>] [--limit <N>] [--page <N>]
+tapd story list [--status <status>] [--owner <owner>] [--iteration-id <id>] [--filter <field=OP<value>>] [--limit <N>] [--page <N>]
```

| 新增标志 | 必需 | 默认值 | 说明 |
|----------|------|--------|------|
| `--filter` | 否 | — | 高级过滤条件（可重复，格式：`field=OP<value>`，支持 LIKE/EQ/CONTAINS 等 OpenAPI 特殊查询语法） |

### tapd bug list

```diff
-tapd bug list [--status <status>] [--priority <priority>] [--severity <severity>] [--limit <N>] [--page <N>]
+tapd bug list [--status <status>] [--priority <priority>] [--severity <severity>] [--filter <field=OP<value>>] [--limit <N>] [--page <N>]
```

| 新增标志 | 必需 | 默认值 | 说明 |
|----------|------|--------|------|
| `--filter` | 否 | — | 高级过滤条件（可重复，格式：`field=OP<value>`，支持 LIKE/EQ/CONTAINS 等 OpenAPI 特殊查询语法） |

### tapd task list

同 story list，新增 `--filter`。

### tapd iteration list

```diff
-tapd iteration list [--status <status>]
+tapd iteration list [--status <status>] [--filter <field=OP<value>>]
```

### 其他 list 命令

以下命令同样新增 `--filter` 标志，格式和行为一致：

- `tapd release list`
- `tapd comment list`
- `tapd timesheet list`
- `tapd attachment list`
- `tapd wiki list`
- `tapd tcase list`
- `tapd category list`
- `tapd custom-field list`
- `tapd workitem-type list`

### 不受影响的命令

- `tapd workspace list` — 响应格式特殊（parseList + category 过滤），且为跨项目查询，无需 filter

## 行为约定

1. `--filter` 可重复使用：`--filter "k1=v1" --filter "k2=v2"`
2. 参数格式为 `key=value`，不含 `=` 的条目静默跳过
3. 不做字段名或操作符校验，非法参数由 TAPD API 返回错误
4. 与已有标志（`--status`、`--owner` 等）可组合使用，参数合并传递
5. `value` 中支持 TAPD 特殊查询操作符：LIKE、EQ、NOT_EQ、LIKE_OR、CONTAINS、CONTAINS_OR、USER_OR，以及时间比较（`>`、`<`、`~`）和多值分隔符（`|`）
