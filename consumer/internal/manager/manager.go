package a2f_manager

import (
	"log"
	"os"
	"path/filepath"

	"example.com/consumer/internal"
	"example.com/consumer/internal/a2f"
	"example.com/consumer/internal/storage"
)

type Manager struct {
	storageManager storage.StorageHandler
	artifactLoader Loader
}

func NewManager(sh storage.StorageHandler) *Manager {
	return &Manager{
		storageManager: sh,
		artifactLoader: NewLoader(),
	}
}

func (m *Manager) SaveArtifact(artifact *a2f.Artifact) error {
	return m.storageManager.SaveArtifact(m.consume(*artifact))
}

func (m *Manager) GetArtifactConf(bId string) (*a2f.ArtifactConf, error) {
	return m.storageManager.GetArtifactConf(bId)
}

func (m *Manager) ListArtifacts() ([]storage.BasicArtifactInfo, error) {
	return m.storageManager.ListArtifacts()
}

// LoadArtifact loads an artifact by its filename
func (m *Manager) LoadArtifact(fname string) (*a2f.Artifact, error) {
	return m.artifactLoader.Load(fname)
}

func (m *Manager) ProcessArtifact(artifact *a2f.Artifact) error {
	return artifact.Process()
}

func (m *Manager) GetArtifactOutput(buildId string) ([]byte, error) {
	basicInfo, err := m.storageManager.GetArtifactMetadata(buildId)
	if err != nil {
		return nil, err
	}

	outputPath := filepath.Join(basicInfo.Location, "output.txt")
	return os.ReadFile(outputPath)
}

func (m *Manager) consume(artifact a2f.Artifact) (a2f.ArtifactConf, []byte, []byte) {

	cnf_loc := artifact.GetConfig().CNFLoc
	atf_loc := filepath.Join(internal.ArtifactDir, artifact.GetConfig().BuildDetails.BuildOutput)

	output, _ := os.ReadFile(filepath.Join(atf_loc, "output.txt"))
	atf, _ := os.ReadFile(filepath.Join(atf_loc, artifact.GetConfig().BinaryObject.Name))

	// Removing
	if err := os.Remove(cnf_loc); err != nil {
		log.Printf("Can't Remove the Artifact Config %s Due to %s\n", cnf_loc, err)
	}

	if err := os.RemoveAll(atf_loc); err != nil {
		log.Printf("Can't Remove the artifacts Dir of %s Due to %s\n", cnf_loc, err)
	}

	return artifact.GetConfig(), output, atf

}
