package a2f

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//////////////////////////////////////////////////////////////////////////
// 						ArtifactConf Structure
//////////////////////////////////////////////////////////////////////////

type Version struct {
	Major int
	Minor int
	Patch int
}

type Metadata struct {
	// Name of the Project
	Name string `json:"name"`
	// Version
	Version Version `json:"version"`
	// time when building the artifact
	BuildTime time.Time `json:"build_time"`
	// Checksum of the A2fBinary Object [ Just to be sure that it's Correct]
	Checksum string `json:"checksum"`
	// ? Where it should be stored .
	Location string `json:"location"`
	// Unique Id
	BuildId string `json:"build_id"`
}

type A2FObject struct {
	// Name of the artifact binary
	Name string `json:"name"`
	// [total space]
	Size int `json:"size"`
}

type BuildInfo struct {
	BuildTime time.Duration `json:"build-time"` // time took to build the artifact

	// Location where the A2f & Output is stored
	BuildOutput string `json:"build-output"`
}

type Step struct {
	// CMD [sh -c] kinda thing
	Cmd string `json:"cmd"`
	// Args
	Input []string `json:"input,omitempty"`
}

type ProcessInfo struct {
	// Group of Execution Steps
	Steps []Step `json:"steps"`
	// ENV Vars loaded berfore exuctuing the Steps
	Env map[string]string `json:"environ,omitempty"`
}

type CookBook struct {
	Build ProcessInfo `json:"build"`
	Exec  ProcessInfo `json:"exec"`
}

type ArtifactConf struct {
	Metadata     Metadata  `json:"metadata"`
	BinaryObject A2FObject `json:"artifact-details"`
	Config       CookBook  `json:"config"`
	BuildDetails BuildInfo `json:"build-details"`
	BuildSuccess bool      `json:"success"`

	// location of the config file [ Config File Name ]
	CNFLoc string `json:"-"`
}

// A2FCnf Impl
func (c ArtifactConf) Name() string         { return c.Metadata.Name }
func (c ArtifactConf) Version() Version     { return c.Metadata.Version }
func (c ArtifactConf) StoreAt() string      { return c.Metadata.Location }
func (c ArtifactConf) Location() string     { return filepath.Join(c.StoreAt(), c.BinaryObject.Name) }
func (c ArtifactConf) ConfigHash() string   { return fileHash(c.CNFLoc) }
func (c ArtifactConf) ArtifactHash() string { return c.Metadata.Checksum }

func fileHash(fpath string) string {
	data, _ := os.ReadFile(fpath)

	hash := sha256.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
