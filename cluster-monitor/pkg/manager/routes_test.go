// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package manager

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gorilla/mux"

	"github.com/couchbaselabs/workbench-prototype/cluster-monitor/pkg/configuration"
)

func createTestFile(t *testing.T, filePath, content string, perms os.FileMode) {
	testStaticFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perms)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testStaticFile.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	if err = testStaticFile.Close(); err != nil {
		t.Fatal(err)
	}
}

func createTestFiles(t *testing.T) string {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(tmpDir, 0o777); err != nil {
		t.Fatal(err)
	}
	createTestFile(t, path.Join(tmpDir, "index.html"), "Index!", 0o774)
	createTestFile(t, path.Join(tmpDir, "test.txt"), "Success!", 0o774)
	createTestFile(t, path.Join(tmpDir, ".hidden.txt"), "Hidden!", 0o774)
	if err := os.Mkdir(path.Join(tmpDir, "directory"), 0o777); err != nil {
		t.Fatal(err)
	}
	return tmpDir
}

func testRequest(router *mux.Router, path string, expectedStatus int, expectedBody string) func(t *testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != expectedStatus {
			t.Fatalf("wrong status code: expected %d, got %d", expectedStatus, status)
		}
		if rr.Body.String() != expectedBody {
			t.Fatalf("unexpected body: got %v, want %v", rr.Body.String(), expectedBody)
		}
	}
}

func TestUIRoutes(t *testing.T) {
	testRoot := createTestFiles(t)
	mgr := &Manager{
		config: &configuration.Config{
			UIRoot: testRoot,
		},
	}
	router := mux.NewRouter()
	ui(router, mgr)

	t.Run("Root", testRequest(router, "/ui/", http.StatusOK, "Index!"))
	t.Run("RootTrailingSlash", testRequest(router, "/ui/", http.StatusOK, "Index!"))

	t.Run("HiddenFile", testRequest(router, "/ui/.hidden.txt", http.StatusForbidden, ""))
	t.Run("Directory", testRequest(router, "/ui/directory", http.StatusForbidden, ""))
	t.Run("DirectoryTrailingSlash", testRequest(router, "/ui/directory/", http.StatusForbidden, ""))

	t.Run("StaticFile", testRequest(router, "/ui/test.txt", http.StatusOK, "Success!"))

	t.Run("NonExistentFile", testRequest(router, "/ui/non-existent", http.StatusOK, "Index!"))
}
