package project

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/builder/internal"
)

func LoadProject(projPath, Id string) (*Project, error) {
	absPath, err := filepath.Abs(projPath)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(absPath); err != nil {
		return nil, err
	}

	configPath := filepath.Join(absPath, "config.json")

	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	var base_envs map[string]string = map[string]string{
		"PROJECT_NAME":  cfg.Project,
		"BUILD_VERSION": cfg.Version,
		"BUILD_ID":      Id,
	}

	return &Project{
		config:   cfg,
		location: projPath,
		envs:     base_envs,
		BuildId:  Id,
	}, nil
}

// project name
func (p *Project) Name() string { return p.config.Project }

// project version
func (p *Project) Version() Version {

	v := strings.Split(p.config.Version, ".")

	major, _ := strconv.Atoi(v[0])
	minor, _ := strconv.Atoi(v[0])
	patch, _ := strconv.Atoi(v[0])

	vr := Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
	return vr
}

// This will return resolved location
// i.e
// if :
//
//	ProjectName = "X"
//	Version = "3.4.5"
//	given :
//		location : $PROJECT_NAME/$BUILD_VERSION
//
// then :
//
//	"X/3.4.5"
func (p *Project) StoreAt() string {
	return os.Expand(p.config.StorageLocation, internal.ResolveENV(p.envs))
}

// artifact location in project
func (p *Project) ArtifactLocation() string {
	return filepath.Join(p.ProjectLocation(), p.config.Artifact)
}

// artifact name
func (p *Project) ArtifactName() string { return p.config.Artifact }

// BaseEnvs::
//
//	 i.e :
//			BUILD_VERSION
//			PROJECT_NAME
func (p *Project) BaseEnvs() map[string]string { return p.envs }

// where the project | codebase is stored
func (p *Project) ProjectLocation() string { return p.location }

// how to build the project
func (p *Project) BuildProcess() ProcessInfo { return p.config.Build }

// how to execute the project
func (p *Project) ExecProcess() ProcessInfo { return p.config.Exec }

// returns the hash of the project
func (p *Project) Hash() []byte { return []byte{} }
