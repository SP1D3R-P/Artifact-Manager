package builder

import "time"

type BuildResult struct {
	ProjectName  string        `json:"build-name"`
	ArtifactPath string        `json:"path"`
	BuildOutput  string        `json:"output"`
	Success      bool          `json:"ok"`
	Error        error         `json:"error,omitempty"`
	BuildTime    time.Duration `json:"build-time"`
}

type Builder interface {
	Build() *BuildResult
}
