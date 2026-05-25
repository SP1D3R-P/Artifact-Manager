package project

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidFormat = errors.New(
		"Invalid Config Format",
	)

	ErrInvalidStorageLocation = errors.New(
		"Invalid Path to storage Location",
	)

	ErrInvalidName = errors.New(
		"Invalid Project Name it must not be Empty.",
	)

	ErrInvalidVersion = errors.New(
		"Invalid Version it must be of format Major.Minor.Patch and it must be of int.",
	)

	ErrInvalidArtifact = errors.New(
		"Invalid Artifact Name. It must not be empty",
	)

	ErrInsufficentBuildSteps = errors.New(
		"Build Steps must not be Less than 1.",
	)

	ErrInsufficentExecSteps = errors.New(
		"Exec Steps must not be Less than 1.",
	)
)

func loadConfig(path string) (*projectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Loaded Config")
	// fmt.Println(string(data))

	var cfg projectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Join(ErrInvalidFormat, err)
	}

	if !verifyName(cfg.Project) {
		return nil, ErrInvalidName
	}

	if !verifyVersion(cfg.Version) {
		return nil, ErrInvalidVersion
	}

	// Here this is just checking if the artifact exists not just empty string
	if !verifyArtifact(cfg.Artifact) {
		return nil, ErrInvalidArtifact
	}

	if !verifyStorageLocation(cfg.StorageLocation) {
		return nil, ErrInvalidStorageLocation
	}

	if len(cfg.Build.Steps) == 0 {
		return nil, ErrInsufficentBuildSteps
	}

	if len(cfg.Exec.Steps) == 0 {
		return nil, ErrInsufficentExecSteps
	}

	return &cfg, nil
}

/*
* Verifier
 */

func verifyName(name string) bool {
	if len(name) == 0 {
		return false
	}

	return true
}

func verifyVersion(v string) bool {
	if len(v) == 0 {
		return false
	}

	versions := strings.Split(v, ".")
	if len(versions) != 3 {
		return false
	}

	for _, m := range versions {
		if _, err := strconv.Atoi(m); err != nil {
			return false
		}
	}

	return true
}

func verifyArtifact(a string) bool        { return len(a) != 0 }
func verifyStorageLocation(a string) bool { return len(a) != 0 }
