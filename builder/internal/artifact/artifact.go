package artifact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"example.com/builder/internal/builder"
	"example.com/builder/internal/project"
)

type Metadata struct {
	Name       string          `json:"name"`
	Version    project.Version `json:"version"`
	BuildTime  time.Time       `json:"build_time"`
	sourcePath string
	Checksum   string `json:"checksum"`
	Location   string `json:"location"`
	BuildId    string `json:"build_id"`
}

type A2FObject struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type BuildInfo struct {
	BuildTime time.Duration `json:"build-time"`

	// Result of the build process
	BuildOutput string `json:"build-output"`
}

type CookBook struct {
	Build project.ProcessInfo `json:"build"`
	Exec  project.ProcessInfo `json:"exec"`
}

type Artifact struct {
	Metadata     Metadata  `json:"metadata"`
	BinaryObject A2FObject `json:"artifact-details"`
	Config       CookBook  `json:"config"`
	BuildDetails BuildInfo `json:"build-details"`
	BuildSuccess bool      `json:"success"`
	filename     string
	result       *builder.BuildResult
	writer       ArtifactWriter
}

func NewArtifact(cnf *project.Project, result *builder.BuildResult) (*Artifact, error) {
	a2fData, err := os.ReadFile(cnf.ArtifactLocation())
	if err != nil {
		return nil, err
	}

	fmt.Println(cnf.StoreAt())
	fmt.Println(filepath.Dir(cnf.StoreAt()))
	fmt.Println(filepath.Base(cnf.StoreAt()))

	return &Artifact{
		Metadata: Metadata{
			Name:       cnf.Name(),
			Version:    cnf.Version(),
			BuildTime:  time.Now(),
			sourcePath: cnf.ProjectLocation(),
			Checksum:   calculateChecksum(a2fData),
			BuildId:    cnf.BuildId,
			Location:   cnf.StoreAt(),
		},
		filename: fmt.Sprintf("%s-%s.json", cnf.Name(), cnf.BuildId),
		Config: CookBook{
			Build: cnf.BuildProcess(),
			Exec:  cnf.ExecProcess(),
		},
		BinaryObject: A2FObject{
			Name: cnf.ArtifactName(),
			Size: len(a2fData),
		},
		BuildDetails: BuildInfo{
			BuildTime: result.BuildTime,
			BuildOutput: filepath.Join(
				filepath.Dir(cnf.StoreAt()),
				fmt.Sprintf("%s-%s", filepath.Base(cnf.StoreAt()), cnf.BuildId),
			),
		},
		writer:       NewBaseArtifactWriter(),
		BuildSuccess: result.Success,
		result:       result,
	}, nil
}

func (a *Artifact) Name() string {
	return a.Metadata.Name
}

func (a *Artifact) Data() []byte {
	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(a)
	if err != nil {
		return nil
	}
	return buffer.Bytes()
}

func (a *Artifact) Filename() string {
	return a.filename
}

func (a *Artifact) SaveTo(dir string) error {
	return a.writer.SaveTo(dir, a)
}
