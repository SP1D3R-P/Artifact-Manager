// Testing the Generic Build
package build4generic

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/builder/internal/project"
)

func TestGenericBuilderBuild(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	testConfig := `{
		"project": "testapp",
		"version": "1.0.0",
		"location": "/tests",
		"artifact": "output.txt",
		"build": {
			"steps": [
				{ "cmd" : "echo test > output.txt" }
			]
		},
		"exec" : {
			"steps" : [
				{ "cmd" : "cat output.txt" }
			]
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// this will be genereated during build time
	outputFile := filepath.Join(tmpDir, "output.txt")

	proj, err := project.LoadProject(tmpDir, "test-build-id")
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	bldr := NewBuilder(proj)
	result := bldr.Build()

	if !result.Success {
		t.Errorf("build should succeed: %v", result.Error)
	}

	if result.ProjectName != "testapp" {
		t.Errorf("expected project name 'testapp', got '%s'", result.ProjectName)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("artifact file not found at %s", outputFile)
	}
}

func TestGenericBuilderBuildWithEnvironment(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	testConfig := `{
		"project": "envtest",
		"version": "1.0.0",
		"location": "a2f",
		"artifact": "result.txt",
		"build": {
			"steps": [
				{ "cmd" : "echo test" }
			],
			"environ": {"CUSTOM_VAR": "custom_value"}
		},
		"exec" : {
			"steps" : [
				{"cmd" : "echo exec"}
			] 
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	resultFile := filepath.Join(tmpDir, "result.txt")
	if err := os.WriteFile(resultFile, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to write artifact: %v", err)
	}

	proj, err := project.LoadProject(tmpDir, "test-build-id")
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	bldr := NewBuilder(proj)
	result := bldr.Build()

	if !result.Success {
		t.Errorf("build with environment should succeed: %v", result.Error)
	}
}

func TestGenericBuilderBuildFailure(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	testConfig := `{
		"project": "failtest",
		"version": "1.0.5",
		"location": "fail/",
		"artifact": "nonexistent.txt",
		"build": {
			"steps": [
				{"cmd" : "exit 1" }
			]
		},
		"exec" : {
			"steps" : [
				{"cmd" : "" }
			]
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	proj, err := project.LoadProject(tmpDir, "test-build-id")
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	bldr := NewBuilder(proj)
	result := bldr.Build()

	if result.Success {
		t.Error("build should fail")
	}

	if result.Error == nil {
		t.Error("error should be set when build fails")
	}
}

func TestGenericBuilderBuildMissingArtifact(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	testConfig := `{
		"project": "noartifact",
		"version": "1.0.0",
		"location": "noartifact",
		"artifact": "missing.txt",
		"build": {
			"steps": [
				{"cmd" : "echo done" }
			]
		},
		"exec" : {
			"steps" : [
				{"cmd" : "" }
			]
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	proj, err := project.LoadProject(tmpDir, "test-build-id")
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	bldr := NewBuilder(proj)
	result := bldr.Build()

	if result.Success {
		t.Error("build should fail when artifact is missing")
	}

	if result.Error == nil {
		t.Error("error should describe missing artifact")
	}

}
