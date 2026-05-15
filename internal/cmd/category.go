// Package cmd 中的 category.go 实现了需求分类管理命令
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var flagCategoryName string

// categoryCmd 是 category 父命令
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "需求分类管理",
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询需求分类列表",
	RunE:  runCategoryList,
}

func init() {
	categoryListCmd.Flags().StringVar(&flagCategoryName, "name", "", "按名称筛选（支持模糊匹配，如 %搜索词%）")
	categoryListCmd.Flags().StringArrayVar(&flagFilter, "filter", nil, "高级过滤条件（可重复，格式：field=OP<value>，支持 LIKE/EQ/CONTAINS 等 OpenAPI 特殊查询语法）")

	categoryCmd.AddCommand(categoryListCmd)
	rootCmd.AddCommand(categoryCmd)
}

func runCategoryList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
	}
	addOptionalParam(params, "name", flagCategoryName)

	categories, err := listWithFilters[model.Category](context.Background(), apiClient, "/story_categories", params, flagFilter, "Category")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items: categories,
		Total: len(categories),
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}
