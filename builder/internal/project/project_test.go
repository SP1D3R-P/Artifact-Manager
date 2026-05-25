package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProject(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	testConfig := `{
		"project": "testproj",
		"version": "1.0.0",
		"location": "output/$PROJECT_NAME/$BUILD_VERSION",
		"artifact": "app.exe",
		"build": {
			"steps": [
				{ "cmd" : "echo test" }
			],
			"environ": {}
		},
		"exec" : {
			"steps" : [
				{"cmd" : "echo 'hello world'"}
			]
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	proj, err := LoadProject(tmpDir, "test-build-id")
	if err != nil {
		t.Fatalf("LoadProject failed: %v", err)
	}

	if proj.Name() != "testproj" {
		t.Errorf("expected name 'testproj', got '%s'", proj.Name())
	}

	var expectedVersion = Version{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
	if proj.Version() != expectedVersion {
		t.Errorf("expected version %v, got '%v'", expectedVersion, proj.Version())
	}

	var expectedLocation = "output/testproj/1.0.0"
	if proj.StoreAt() != expectedLocation {
		t.Errorf("expected location %s , got '%s'", expectedLocation, proj.StoreAt())
	}
}
