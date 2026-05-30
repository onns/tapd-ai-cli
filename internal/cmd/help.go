// Package cmd 中的 spec.go 实现了紧凑参考卡生成逻辑，供根命令 --help 使用
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// groupPriority 定义命令分组的显示优先级，数字越小越靠前
// 未列出的分组默认优先级为 100，按字母序排列
var groupPriority = map[string]int{
	"url":     1,
	"story":   2,
	"comment": 3,
	"task":    4,
	"bug":     5,
}

// coreCommands 定义 v0.7.0 版本中已有的核心命令集合
// key 为 group 名（顶层命令），value 为该 group 下的子命令名集合
// 叶子命令（如 url）使用自身名称作为子命令名
var coreCommands = map[string]map[string]bool{
	"url":          {"url": true},
	"story":        {"list": true, "show": true, "create": true, "update": true, "count": true, "todo": true},
	"comment":      {"list": true, "add": true, "update": true, "count": true},
	"task":         {"list": true, "show": true, "create": true, "update": true, "count": true, "todo": true},
	"bug":          {"list": true, "show": true, "create": true, "update": true, "count": true, "todo": true},
	"wiki":         {"list": true, "show": true, "create": true, "update": true},
	"iteration":    {"list": true, "create": true, "update": true, "count": true},
	"tcase":        {"list": true, "create": true, "batch-create": true},
	"timesheet":    {"list": true, "add": true, "update": true},
	"workflow":     {"transitions": true, "status-map": true, "last-steps": true},
	"relation":     {"bugs": true, "create": true},
	"auth":         {"login": true},
	"workspace":    {"list": true, "switch": true, "info": true},
	"attachment":   {"list": true},
	"image":        {"get": true},
	"category":     {"list": true},
	"custom-field": {"list": true},
	"story-field":  {"info": true, "label": true},
	"workitem-type": {"list": true},
	"release":      {"list": true},
	"skill":        {"init": true},
	"qiwei":        {"send": true},
	"commit-msg":   {"get": true},
}

// specLine 表示一条命令参考行
type specLine struct {
	group    string // 命令所属分组（第一级子命令名）
	leafName string // 叶子命令名（路径最后一段）
	text     string // 完整的命令参考文本
}

// buildSpecLines 遍历命令树，为每个叶子命令生成参考行，按核心/高级分组并排序
func buildSpecLines(root *cobra.Command) (coreLines, advancedLines []specLine) {
	var allLines []specLine
	walkSpecCommands(root, "", "", &allLines)
	for _, l := range allLines {
		if isCoreCommand(l.group, l.leafName) {
			coreLines = append(coreLines, l)
		} else {
			advancedLines = append(advancedLines, l)
		}
	}
	sortSpecLines(coreLines)
	sortSpecLines(advancedLines)
	return
}

// isCoreCommand 判断命令是否属于 v0.7.0 核心命令集
func isCoreCommand(group, leafName string) bool {
	subs, ok := coreCommands[group]
	if !ok {
		return false
	}
	return subs[leafName]
}

// sortSpecLines 按 groupPriority 对参考行排序，优先级相同的按分组名字母序
func sortSpecLines(lines []specLine) {
	sort.SliceStable(lines, func(i, j int) bool {
		pi := getGroupPriority(lines[i].group)
		pj := getGroupPriority(lines[j].group)
		if pi != pj {
			return pi < pj
		}
		return lines[i].group < lines[j].group
	})
}

// getGroupPriority 返回分组的显示优先级，未配置的分组返回默认值 100
func getGroupPriority(group string) int {
	if p, ok := groupPriority[group]; ok {
		return p
	}
	return 100
}

// walkSpecCommands 递归遍历命令树，收集叶子命令的参考行
func walkSpecCommands(cmd *cobra.Command, prefix string, group string, lines *[]specLine) {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}

		fullPath := child.Name()
		if prefix != "" {
			fullPath = prefix + " " + child.Name()
		}

		// 确定分组名：取第一级子命令名
		currentGroup := group
		if currentGroup == "" {
			currentGroup = child.Name()
		}

		if child.HasSubCommands() {
			walkSpecCommands(child, fullPath, currentGroup, lines)
		} else {
			line := commandToLine(child, fullPath)
			// leafName 为相对于 group 的子命令路径
			// 例如 group="story", fullPath="story list" → leafName="list"
			// 例如 group="story", fullPath="story link list" → leafName="link list"
			leafName := fullPath
			if currentGroup != "" && strings.HasPrefix(fullPath, currentGroup+" ") {
				leafName = fullPath[len(currentGroup)+1:]
			}
			*lines = append(*lines, specLine{group: currentGroup, leafName: leafName, text: line})
		}
	}
}

// commandToLine 将 Cobra 命令转换为一行紧凑参考文本
func commandToLine(cmd *cobra.Command, path string) string {
	var b strings.Builder
	b.WriteString("tapd ")
	b.WriteString(path)

	// 添加位置参数
	if argName := extractArgName(cmd.Use); argName != "" {
		b.WriteString(" <")
		b.WriteString(argName)
		b.WriteString(">")
	}

	// 检测是否同时有 --description 和 --file（富文本输入命令）
	hasDescription := cmd.Flags().Lookup("description") != nil
	hasFile := cmd.Flags().Lookup("file") != nil
	richTextWritten := false

	// 收集标志，区分必填和可选
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		// 跳过全局认证标志（spec 不需要展示）
		if isGlobalAuthFlag(f.Name) {
			return
		}
		// 跳过全局展示标志（已在 header 中展示）
		if isGlobalDisplayFlag(f.Name) {
			return
		}
		// 富文本输入：合并 --description / --file / stdin 为一个组合提示
		if hasDescription && hasFile && (f.Name == "description" || f.Name == "file") {
			if !richTextWritten {
				b.WriteString(" [--description=<text>|--file=<path>|stdin]")
				richTextWritten = true
			}
			return
		}
		b.WriteString(" ")
		b.WriteString(formatFlag(f))
	})

	// 添加描述注释
	if cmd.Short != "" {
		b.WriteString("  # ")
		b.WriteString(cmd.Short)
	}

	return b.String()
}

// formatFlag 将一个 flag 格式化为紧凑文本
// 必填标志：--flag=<val>
// 可选带枚举：[--flag=<a|b|c>]
// 可选带默认值：[--flag=default]
// 可选无默认值：[--flag]
func formatFlag(f *pflag.Flag) string {
	hint := extractEnumHint(f.Usage)
	if isFlagRequired(f) {
		if hint != "" {
			return "--" + f.Name + "=<" + hint + ">"
		}
		return "--" + f.Name + "=<" + f.Name + ">"
	}
	if hint != "" {
		return "[--" + f.Name + "=<" + hint + ">]"
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
		return "[--" + f.Name + "=" + f.DefValue + "]"
	}
	return "[--" + f.Name + "]"
}

// extractEnumHint 从 Usage 文本的全角括号 （）中提取枚举提示
// 例如 "优先级（urgent/high/medium/low）" -> "urgent/high/medium/low"
// 若括号内含必填说明（，必需/，必填）则去掉后缀
func extractEnumHint(usage string) string {
	start := strings.Index(usage, "（")
	end := strings.LastIndex(usage, "）")
	if start < 0 || end <= start {
		return ""
	}
	content := usage[start+len("（") : end]
	content = strings.TrimSuffix(content, "，必需")
	content = strings.TrimSuffix(content, "，必填")
	return content
}

// isFlagRequired 判断标志是否为必填（通过检测 Usage 文本中的关键字）
func isFlagRequired(f *pflag.Flag) bool {
	usage := f.Usage
	return strings.Contains(usage, "必需") || strings.Contains(usage, "必填")
}

// isGlobalAuthFlag 判断是否为全局认证标志
func isGlobalAuthFlag(name string) bool {
	switch name {
	case "access-token", "api-user", "api-password":
		return true
	default:
		return false
	}
}

// isGlobalDisplayFlag 判断是否为全局展示标志（header 中已展示）
func isGlobalDisplayFlag(name string) bool {
	switch name {
	case "workspace-id", "pretty", "json", "no-comments":
		return true
	default:
		return false
	}
}

// printSpecOutput 输出完整的参考卡文本，分为核心命令和高级命令两个区域
func printSpecOutput(w *os.File, root *cobra.Command, coreLines, advancedLines []specLine) {
	// 标题行
	fmt.Fprintf(w, "tapd - %s\n", root.Short)
	fmt.Fprintln(w, "Global: [--workspace-id=<id>] [--json（详情提取字段时用，默认勿加）] [--pretty（人类阅读用，AI勿加）] [--no-comments]")

	// 输出核心命令
	printGroupedLines(w, coreLines)

	// 输出高级命令
	if len(advancedLines) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "# ─── 高级命令 ───")
		printGroupedLines(w, advancedLines)
	}
}

// printGroupedLines 按分组输出命令参考行
func printGroupedLines(w *os.File, lines []specLine) {
	lastGroup := ""
	for _, l := range lines {
		if l.group != lastGroup {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "# %s\n", l.group)
			lastGroup = l.group
		}
		fmt.Fprintln(w, l.text)
	}
}

// extractArgName 从 Use 字段提取位置参数名（如 "show <story_id>" -> "story_id"）
func extractArgName(use string) string {
	start := -1
	for i, c := range use {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start >= 0 {
			return use[start:i]
		}
	}
	return ""
}
