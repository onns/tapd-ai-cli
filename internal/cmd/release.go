// Package cmd 中的 release.go 实现了发布计划管理命令
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

// releaseCmd 是 release 父命令
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "发布计划管理",
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询发布计划列表",
	RunE:  runReleaseList,
}

func init() {
	releaseListCmd.Flags().StringArrayVar(&flagFilter, "filter", nil, "高级过滤条件（可重复，格式：field=OP<value>，支持 LIKE/EQ/CONTAINS 等 OpenAPI 特殊查询语法）")

	releaseCmd.AddCommand(releaseListCmd)
	rootCmd.AddCommand(releaseCmd)
}

func runReleaseList(cmd *cobra.Command, args []string) error {
	req := &model.WorkspaceIDRequest{
		WorkspaceID: flagWorkspaceID,
	}

	releases, err := listWithFilters[model.Release](context.Background(), apiClient, "/releases", req.ToParams(), flagFilter, "Release")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, releases, !flagPretty)
}
