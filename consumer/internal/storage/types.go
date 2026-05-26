package storage

import (
	"example.com/consumer/internal/a2f"
)

type BasicArtifactInfo struct {
	BuildId      string
	Name         string
	Version      a2f.Version
	ArtifactName string

	Success       bool
	BuildTime     string
	BuildDuration string

	// Checksums ::
	// 		not sure if i should provide this or not
	ConfigChecksum   string
	ArtifactChecksum string
}

type StorageHandler interface {
	GetArtifactMetadata(buildId string) (*ArtifactMetadata, error)
	// SaveArtifact saves an artifact to the central repository
	SaveArtifact(cnf a2f.ArtifactConf, output []byte, atf []byte) error
	// GetArtifactConf retrieves an artifact config by its build ID
	GetArtifactConf(bId string) (*a2f.ArtifactConf, error)
	// ListArtifacts lists all artifacts in the central repository
	ListArtifacts() ([]BasicArtifactInfo, error)
}
