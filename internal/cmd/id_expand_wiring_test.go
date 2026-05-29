// Package cmd 中的 id_expand_wiring_test.go 验证 show / update 命令在调用 SDK 前
// 已经把短 ID 展开为 TAPD 长 ID（即把短号 28841 + workspace 12345 → 1112345000028841 再发请求）。
// 这些测试拦截 HTTP 请求并断言出栈到 SDK 的 id 参数已经是展开后的长 ID。
package cmd

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// captureIDs 记录每次非 /comments 请求中的 id：
//   - GET（show 系列）从 URL query 取
//   - POST（update 系列）从表单体取
// /comments 请求单独记录 entry_id，便于断言 printComments 也用了展开后的 ID
type captureIDs struct {
	mainID    string // GetBug/GetStory/GetTask 或 UpdateXxx 收到的 id
	commentID string // ListComments 收到的 entry_id
}

// captureHandler 返回一个测试服务器 handler，写入 caps 同时回一份合法响应让 handler 流程跑完
func captureHandler(t *testing.T, entity string, caps *captureIDs) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/comments") {
			caps.commentID = r.URL.Query().Get("entry_id")
			w.Write([]byte(`{"status":1,"data":[]}`))
			return
		}
		switch r.Method {
		case http.MethodGet:
			caps.mainID = r.URL.Query().Get("id")
			// GET (show/list) 返回数组 shape
			w.Write([]byte(`{"status":1,"data":[{"` + entity + `":{"id":"placeholder","title":"x","name":"x"}}]}`))
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			form, _ := url.ParseQuery(string(body))
			caps.mainID = form.Get("id")
			// POST (update) 返回单对象 shape
			w.Write([]byte(`{"status":1,"data":{"` + entity + `":{"id":"placeholder","title":"x","name":"x"}}}`))
		}
	}
}

func TestRunBugShow_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Bug", &caps))
	defer cleanup()
	flagJSON = true
	flagNoComments = false // 让 printComments 被实际触发，以验证它也用了展开后的 ID

	restore, reader := captureStdout(t)
	if err := runBugShow(nil, []string{"28841"}); err != nil {
		t.Fatalf("runBugShow returned error: %v", err)
	}
	restore()
	drainReader(reader)

	// setupMockServer 固定 workspaceID = "12345"，按 expandShortID 规则：
	// "11" + "12345" + 左补零到 9 位的 "28841" = "1112345000028841"
	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("GetBug received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
	if caps.commentID != wantID {
		t.Errorf("printComments received entry_id=%q, want %q (短号未被展开)", caps.commentID, wantID)
	}
}

func TestRunStoryShow_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Story", &caps))
	defer cleanup()
	flagJSON = true
	flagNoComments = false

	restore, reader := captureStdout(t)
	if err := runStoryShow(nil, []string{"28841"}); err != nil {
		t.Fatalf("runStoryShow returned error: %v", err)
	}
	restore()
	drainReader(reader)

	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("GetStory received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
	if caps.commentID != wantID {
		t.Errorf("printComments received entry_id=%q, want %q (短号未被展开)", caps.commentID, wantID)
	}
}

func TestRunTaskShow_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Task", &caps))
	defer cleanup()
	flagJSON = true
	flagNoComments = false

	restore, reader := captureStdout(t)
	if err := runTaskShow(nil, []string{"28841"}); err != nil {
		t.Fatalf("runTaskShow returned error: %v", err)
	}
	restore()
	drainReader(reader)

	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("GetTask received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
	if caps.commentID != wantID {
		t.Errorf("printComments received entry_id=%q, want %q (短号未被展开)", caps.commentID, wantID)
	}
}

func TestRunBugUpdate_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Bug", &caps))
	defer cleanup()
	flagJSON = true
	flagTitle = "x" // 给 update 一个非空字段以确保 SDK 真的发请求

	restore, reader := captureStdout(t)
	if err := runBugUpdate(nil, []string{"28841"}); err != nil {
		t.Fatalf("runBugUpdate returned error: %v", err)
	}
	restore()
	drainReader(reader)

	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("UpdateBug received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
}

func TestRunStoryUpdate_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Story", &caps))
	defer cleanup()
	flagJSON = true
	flagName = "x"

	restore, reader := captureStdout(t)
	if err := runStoryUpdate(nil, []string{"28841"}); err != nil {
		t.Fatalf("runStoryUpdate returned error: %v", err)
	}
	restore()
	drainReader(reader)

	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("UpdateStory received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
}

func TestRunTaskUpdate_ExpandsShortID(t *testing.T) {
	resetFlags()
	var caps captureIDs
	_, cleanup := setupMockServer(t, captureHandler(t, "Task", &caps))
	defer cleanup()
	flagJSON = true
	flagName = "x"

	restore, reader := captureStdout(t)
	if err := runTaskUpdate(nil, []string{"28841"}); err != nil {
		t.Fatalf("runTaskUpdate returned error: %v", err)
	}
	restore()
	drainReader(reader)

	const wantID = "1112345000028841"
	if caps.mainID != wantID {
		t.Errorf("UpdateTask received id=%q, want %q (短号未被展开)", caps.mainID, wantID)
	}
}
