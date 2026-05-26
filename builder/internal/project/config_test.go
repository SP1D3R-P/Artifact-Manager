// =================================
//
//	Testing the Config
//
// =================================
//
// this Contains :
//
// # TestLoadConfig :
//		Loading the a valid Config file and checkng if everyting is loaded properly
//
// # TestLoadConfigNotFound :
//		Trying to load a config file where it's not present
//
// # TestLoadConfigInvalidJSON :
// 		Loading a config file with invalid format

package project

import (
	"os"
	"path/filepath"
	"testing"
)

// Testing all the properties of the Config is loaded correctly
func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := `{
		"project": "test",
		"version": "1.0.0",
		"location": "test/$PROJECT_NAME/$BUILD_VERSION",
		"artifact": "main.exe",
		"build": {
			"steps": [
				{"cmd" : "echo hello"}
			],
			"environ": {}
		},
		"exec" : {
			"steps" : [
				{"cmd" : "echo hello"}
			]
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), os.ModePerm); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Project != "test" {
		t.Errorf("expected project 'test', got '%s'", cfg.Project)
	}

	if cfg.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", cfg.Version)
	}

	if len(cfg.Build.Steps) != 1 {
		t.Errorf("expected 1 build step, got %d", len(cfg.Build.Steps))
	}

	if cfg.StorageLocation != "test/$PROJECT_NAME/$BUILD_VERSION" {
		t.Errorf("expected location with template, got '%s'", cfg.StorageLocation)
	}
}

// Loading Config Not present
func TestLoadConfigNotFound(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing config, got nil")
	}
}

// Invalid Structure
func TestLoadConfigInvalidJSON(t *testing.T) {
	// TODO :: Use more invalid format
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := `{invalid json}`

	if err := os.WriteFile(configPath, []byte(testConfig), os.ModePerm); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := loadConfig(configPath)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
