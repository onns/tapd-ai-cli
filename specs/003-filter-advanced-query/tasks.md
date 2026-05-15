# 任务: 所有 list 命令增加 --filter 高级过滤

**输入**: 来自 `/specs/003-filter-advanced-query/` 的设计文档
**前置条件**: spec.md（必需）, contracts/filter-contract.md（必需）

## 格式: `[ID] [Story] 描述`
- **[Story]**: US1=按自定义字段过滤, US2=时间范围组合查询, US3=所有 list 命令统一支持

---

## 阶段 1: SDK 改动

- [x] T001 在 `tapd-sdk-go/client.go` 新增公开 `DoGet` 方法
- [x] T002 在 `tapd-sdk-go/parse.go` 新增公开 `ParseList` 泛型函数
- [x] T003 在主 `go.mod` 添加 `replace github.com/studyzy/tapd-sdk-go => ../tapd-sdk-go`

---

## 阶段 2: CLI 核心 helper

- [x] T004 在 `internal/cmd/root.go` 新增 `flagFilter` 全局变量、`listWithFilters[T]` 泛型 helper
- [x] T005 编译验证 `go build ./...`

---

## 阶段 3: 各命令适配

- [x] T006 `internal/cmd/story.go` — 注册 `--filter`，改用 `listWithFilters` [US1]
- [x] T007 `internal/cmd/bug.go` — 注册 `--filter`，改用 `listWithFilters` [US2]
- [x] T008 `internal/cmd/task.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T009 `internal/cmd/iteration.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T010 `internal/cmd/release.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T011 `internal/cmd/comment.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T012 `internal/cmd/timesheet.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T013 `internal/cmd/attachment.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T014 `internal/cmd/wiki.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T015 `internal/cmd/tcase.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T016 `internal/cmd/category.go` — 注册 `--filter`，改用 `listWithFilters` [US3]
- [x] T017 `internal/cmd/custom_field.go` — 注册 `--filter`，改用 `listWithFilters`（custom-field + workitem-type） [US3]

---

## 阶段 4: 文档

- [x] T018 `internal/cmd/skill_template.md` — 新增"高级过滤（--filter）"章节
- [x] T019 `README.md` — 新增"高级过滤（--filter）"章节
- [x] T020 `CODEBUDDY.md` — 核心设计原则新增第 7 点
- [x] T021 `specs/003-filter-advanced-query/` — 新增 spec.md、contracts/filter-contract.md、tasks.md
