// Package cmd 中的 new_commands2_test.go 为剩余新增命令提供 mock server 单元测试
package cmd

import (
	"net/http"
	"strings"
	"testing"
)

// newCmds2Handler 统一处理第二批新增命令的 mock API 响应
func newCmds2Handler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	path := r.URL.Path

	switch {
	// ===== batch_update =====
	case strings.HasSuffix(path, "/stories/batch_update_story"):
		w.Write([]byte(`{"status":1,"data":{"msg":"ok"}}`))
	case strings.HasSuffix(path, "/bugs/batch_update_bug"):
		w.Write([]byte(`{"status":1,"data":""}`))
	case strings.HasSuffix(path, "/tasks/batch_update_task"):
		w.Write([]byte(`{"status":1,"data":{"msg":"ok"}}`))

	// ===== removed =====
	case strings.HasSuffix(path, "/stories/get_removed_stories"):
		w.Write([]byte(`{"status":1,"data":[{"RemovedStory":{"id":"1","name":"removed"}}]}`))
	case strings.HasSuffix(path, "/bugs/get_removed_bugs"):
		w.Write([]byte(`{"status":1,"data":[{"RemovedBug":{"id":"1","title":"removed"}}]}`))
	case strings.HasSuffix(path, "/tasks/get_removed_tasks"):
		w.Write([]byte(`{"status":1,"data":[{"RemovedTask":{"id":"1","name":"removed"}}]}`))

	// ===== release =====
	case strings.HasSuffix(path, "/new_releases/count"):
		w.Write([]byte(`{"status":1,"data":{"count":3}}`))
	case strings.Contains(path, "/new_releases") || strings.HasSuffix(path, "/releases"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"Release":{"id":"1","name":"v1.0"}}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"Release":{"id":"1","name":"v1.0"}}]}`))
		}

	// ===== workflow =====
	case strings.HasSuffix(path, "/workflows/first_step"):
		w.Write([]byte(`{"status":1,"data":{"open":"新建"}}`))
	case strings.HasSuffix(path, "/workflows"):
		w.Write([]byte(`{"status":1,"data":[{"Workflow":{"id":"1","name":"默认工作流","system_name":"story"}}]}`))

	// ===== source =====
	case strings.HasSuffix(path, "/code_commit_infos"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"id":"1001","workspace_id":"12345","message":"fix bug","author":"admin"}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"id":"1001","workspace_id":"12345","message":"fix bug","author":"admin"}]}`))
		}
	case strings.HasSuffix(path, "/code_commit_objects/workitems"):
		w.Write([]byte(`{"status":1,"data":{}}`))

	// ===== category =====
	case strings.HasSuffix(path, "/story_categories/count"):
		w.Write([]byte(`{"status":1,"data":{"count":10}}`))
	case strings.HasSuffix(path, "/story_categories"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"Category":{"id":"1","name":"cat1"}}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"Category":{"id":"1","name":"cat1"}}]}`))
		}

	// ===== version_mgmt =====
	case strings.HasSuffix(path, "/versions/count"):
		w.Write([]byte(`{"status":1,"data":{"count":2}}`))
	case strings.HasSuffix(path, "/versions"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"Version":{"id":"1","name":"v1","workspace_id":"12345"}}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"Version":{"id":"1","name":"v1"}}]}`))
		}

	// ===== baseline =====
	case strings.HasSuffix(path, "/baselines/count"):
		w.Write([]byte(`{"status":1,"data":{"count":4}}`))
	case strings.HasSuffix(path, "/baselines"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"Baseline":{"id":"1","name":"bl1","workspace_id":"12345"}}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"Baseline":{"id":"1","name":"bl1"}}]}`))
		}

	// ===== feature =====
	case strings.HasSuffix(path, "/features/count"):
		w.Write([]byte(`{"status":1,"data":{"count":6}}`))
	case strings.HasSuffix(path, "/features"):
		if r.Method == http.MethodPost {
			w.Write([]byte(`{"status":1,"data":{"Feature":{"id":"1","name":"feat1","workspace_id":"12345"}}}`))
		} else {
			w.Write([]byte(`{"status":1,"data":[{"Feature":{"id":"1","name":"feat1"}}]}`))
		}

	// ===== iteration lock/unlock =====
	case strings.HasSuffix(path, "/iterations/lock_iteration"):
		w.Write([]byte(`{"status":1,"info":"success","data":""}`))
	case strings.HasSuffix(path, "/iterations/unlock_iteration"):
		w.Write([]byte(`{"status":1,"info":"success","data":""}`))
	case strings.HasSuffix(path, "/iterations/count"):
		w.Write([]byte(`{"status":1,"data":{"count":1}}`))

	// ===== report =====
	case strings.HasSuffix(path, "/workspace_reports"):
		w.Write([]byte(`{"status":1,"data":[{"WorkspaceReport":{"id":"1","name":"周报"}}]}`))

	// ===== wiki =====
	case strings.HasSuffix(path, "/tapd_wikis/count"):
		w.Write([]byte(`{"status":1,"data":{"count":7}}`))

	// ===== launch logs =====
	case strings.HasSuffix(path, "/launch_forms/get_activity_logs"):
		w.Write([]byte(`{"status":1,"data":[{"LaunchFormActivityLog":{"id":"1"}}]}`))
	case strings.HasSuffix(path, "/launch_forms/count"):
		w.Write([]byte(`{"status":1,"data":{"count":0}}`))

	// ===== story_time_relation =====
	case strings.HasSuffix(path, "/stories/get_time_relative_stories"):
		w.Write([]byte(`{"status":1,"data":[{"WorkitemTimeRelation":{"id":"1"}}]}`))
	case strings.HasSuffix(path, "/stories/delete_time_relations"):
		w.Write([]byte(`{"status":1,"data":{"num":1}}`))

	// ===== story_template =====
	case strings.HasSuffix(path, "/stories/template_list"):
		w.Write([]byte(`{"status":1,"data":[{"WorkitemTemplate":{"id":"1","name":"tmpl"}}]}`))

	default:
		w.Write([]byte(`{"status":1,"data":{}}`))
	}
}

// ===================== batch_update 测试 =====================

func TestNew2RunStoryBatchUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagBatchIDs = "10001,10002"
	flagStatus = "done"

	restore, reader := captureStdout(t)
	err := runStoryBatchUpdate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryBatchUpdate failed: %v", err)
	}
}

func TestNew2RunBugBatchUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagBatchIDs = "20001,20002"
	flagStatus = "fixed"

	restore, reader := captureStdout(t)
	err := runBugBatchUpdate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugBatchUpdate failed: %v", err)
	}
}

func TestNew2RunTaskBatchUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagBatchIDs = "30001,30002"
	flagStatus = "done"

	restore, reader := captureStdout(t)
	err := runTaskBatchUpdate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskBatchUpdate failed: %v", err)
	}
}

// ===================== removed 测试 =====================

func TestNew2RunStoryRemoved(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runStoryRemoved(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryRemoved failed: %v", err)
	}
}

func TestNew2RunBugRemoved(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runBugRemoved(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugRemoved failed: %v", err)
	}
}

func TestNew2RunTaskRemoved(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runTaskRemoved(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskRemoved failed: %v", err)
	}
}

// ===================== release 测试 =====================

func TestNew2RunReleaseCreate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "v1.0"
	flagStartDate = "2026-01-01"
	flagEndDate = "2026-01-31"

	restore, reader := captureStdout(t)
	err := runReleaseCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runReleaseCreate failed: %v", err)
	}
}

func TestNew2RunReleaseUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "v2.0"
	flagStatus = "active"

	restore, reader := captureStdout(t)
	err := runReleaseUpdate(nil, []string{"1001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runReleaseUpdate failed: %v", err)
	}
}

func TestNew2RunReleaseCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runReleaseCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runReleaseCount failed: %v", err)
	}
}

// ===================== workflow 测试 =====================

func TestNew2RunWorkflowFirstStep(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagSystem = "story"

	restore, reader := captureStdout(t)
	err := runWorkflowFirstStep(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkflowFirstStep failed: %v", err)
	}
}

func TestNew2RunWorkflowList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagSystem = "story"

	restore, reader := captureStdout(t)
	err := runWorkflowList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkflowList failed: %v", err)
	}
}

// ===================== source 测试 =====================

func TestNew2RunSourceAdd(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagMessage = "feat: add login"
	flagWebURL = "https://git.example.com/repo"
	flagRef = "refs/heads/main"

	restore, reader := captureStdout(t)
	err := runSourceAdd(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runSourceAdd failed: %v", err)
	}
}

func TestNew2RunSourceList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagObjectID = "10001"
	flagEntityType = "story"

	restore, reader := captureStdout(t)
	err := runSourceList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runSourceList failed: %v", err)
	}
}

func TestNew2RunSourceObjects(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagCommitID = "abc123"
	flagEntityType = "story"

	restore, reader := captureStdout(t)
	err := runSourceObjects(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runSourceObjects failed: %v", err)
	}
}

// ===================== category 测试 =====================

func TestNew2RunCategoryCreate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "NewCat"
	flagParentID = "0"

	restore, reader := captureStdout(t)
	err := runCategoryCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runCategoryCreate failed: %v", err)
	}
}

func TestNew2RunCategoryUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "UpdatedCat"

	restore, reader := captureStdout(t)
	err := runCategoryUpdate(nil, []string{"1001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runCategoryUpdate failed: %v", err)
	}
}

func TestNew2RunCategoryCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runCategoryCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runCategoryCount failed: %v", err)
	}
}

// ===================== version_mgmt 测试 =====================

func TestNew2RunAppVersionList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runAppVersionList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runAppVersionList failed: %v", err)
	}
}

func TestNew2RunAppVersionCreate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "v1.0"
	flagCreator = "admin"

	restore, reader := captureStdout(t)
	err := runAppVersionCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runAppVersionCreate failed: %v", err)
	}
}

func TestNew2RunAppVersionUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagVersionModifier = "admin"
	flagName = "v2.0"

	restore, reader := captureStdout(t)
	err := runAppVersionUpdate(nil, []string{"1001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runAppVersionUpdate failed: %v", err)
	}
}

func TestNew2RunAppVersionCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runAppVersionCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runAppVersionCount failed: %v", err)
	}
}

// ===================== baseline 测试 =====================

func TestNew2RunBaselineList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runBaselineList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBaselineList failed: %v", err)
	}
}

func TestNew2RunBaselineCreate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "baseline1"
	flagVersionID = "v001"

	restore, reader := captureStdout(t)
	err := runBaselineCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBaselineCreate failed: %v", err)
	}
}

func TestNew2RunBaselineUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "baseline-updated"

	restore, reader := captureStdout(t)
	err := runBaselineUpdate(nil, []string{"1001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBaselineUpdate failed: %v", err)
	}
}

func TestNew2RunBaselineCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runBaselineCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBaselineCount failed: %v", err)
	}
}

// ===================== feature 测试 =====================

func TestNew2RunFeatureList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runFeatureList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runFeatureList failed: %v", err)
	}
}

func TestNew2RunFeatureCreate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "NewFeature"

	restore, reader := captureStdout(t)
	err := runFeatureCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runFeatureCreate failed: %v", err)
	}
}

func TestNew2RunFeatureUpdate(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagName = "UpdatedFeature"
	flagOwner = "dev1"

	restore, reader := captureStdout(t)
	err := runFeatureUpdate(nil, []string{"1001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runFeatureUpdate failed: %v", err)
	}
}

func TestNew2RunFeatureCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runFeatureCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runFeatureCount failed: %v", err)
	}
}

// ===================== iteration lock/unlock 测试 =====================

func TestNew2RunIterationLock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLockTypes = "__ALL_STORY__"

	restore, reader := captureStdout(t)
	err := runIterationLock(nil, []string{"40001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationLock failed: %v", err)
	}
}

func TestNew2RunIterationUnlock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runIterationUnlock(nil, []string{"40001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationUnlock failed: %v", err)
	}
}

// ===================== report 测试 =====================

func TestNew2RunReportList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLimit = 30
	flagPage = 1

	restore, reader := captureStdout(t)
	err := runReportList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runReportList failed: %v", err)
	}
}

// ===================== wiki count 测试 =====================

func TestNew2RunWikiCount(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runWikiCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWikiCount failed: %v", err)
	}
}

// ===================== launch logs 测试 =====================

func TestNew2RunLaunchLogs(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagLaunchID = "50001"

	restore, reader := captureStdout(t)
	err := runLaunchLogs(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runLaunchLogs failed: %v", err)
	}
}

// ===================== story_time_relation 测试 =====================

func TestNew2RunStoryTimeRelationList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagTimeRelStoryID = "10001"

	restore, reader := captureStdout(t)
	err := runStoryTimeRelationList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryTimeRelationList failed: %v", err)
	}
}

func TestNew2RunStoryTimeRelationDelete(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	flagTimeRelStoryID = "10001"
	flagTimeRelRelatedID = "10002"
	flagCurrentUser = "admin"

	restore, reader := captureStdout(t)
	err := runStoryTimeRelationDelete(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryTimeRelationDelete failed: %v", err)
	}
}

// ===================== story_template 测试 =====================

func TestNew2RunStoryTemplateList(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, newCmds2Handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runStoryTemplateList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryTemplateList failed: %v", err)
	}
}
