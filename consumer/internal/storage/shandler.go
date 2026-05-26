package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"example.com/consumer/internal"
	"example.com/consumer/internal/a2f"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ArtifactMetadata stores artifact details in the database
type ArtifactMetadata struct {
	BuildID string `gorm:"primaryKey"`
	Name    string

	// version
	VersionMajor int
	VersionMinor int
	VersionPatch int

	BuildAt            time.Time
	BuildTime          time.Duration
	ConfigFileLocation string // location of the config file in the filesystem
	Checksum           string // Checksum of the Config
	Success            bool   // Whether the artifact was successfully processed

	Location string // where the code is stored in the filesystem
	// Contents -
	// {Location}/{BinaryObject.Name}
	// {Location}/output.txt

	ArtifactName string // name of the artifact binary
	ArtifactHash string // hash of the artifact binary
}

type basicHandler struct {
	db          *gorm.DB
	storagePath string
}

func NewBasicHandler() StorageHandler {
	handler := &basicHandler{}
	if err := handler.initDB(); err != nil {
		log.Printf("ERROR :: Error initializing database: %v\n", err)
	}
	return handler
}

func (s *basicHandler) initDB() error {

	storagePath := internal.StorageFSDir
	s.storagePath = storagePath

	dbPath := filepath.Join(storagePath, "artifacts.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to open database: %w", err)
	}

	s.db = db

	if err := s.db.AutoMigrate(&ArtifactMetadata{}); err != nil {
		return fmt.Errorf("Failed to migrate database: %w", err)
	}

	log.Printf("INFO :: Database initialized successfully at %s\n", dbPath)
	return nil
}

// TODOs :: [NOT IMPORTANT FOR WORKING PROTOTYPE]
// 1. Undo All the writing if any of the step fails
// 2. But Also add to database that it failed
func (s *basicHandler) SaveArtifact(cnf a2f.ArtifactConf, output, atf []byte) error {
	log.Printf("INFO :: Saving the Artifact %s\n", cnf.Name())

	// where the artifact will be stored in the filesystem
	// making this
	//  	artifact
	//		configs
	// seperate to avoid confusion and to make it easier to manage

	artifactDir := filepath.Join(s.storagePath, "artifacts", cnf.StoreAt())
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		return fmt.Errorf("failed to create artifact directory: %w", err)
	}

	configDir := filepath.Join(s.storagePath, "configs", cnf.Metadata.BuildId)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create configs directory: %w", err)
	}

	// main artifact
	artifactPath := filepath.Join(artifactDir, cnf.BinaryObject.Name)
	if err := os.WriteFile(artifactPath, atf, 0644); err != nil {
		return fmt.Errorf("failed to save artifact file: %w", err)
	}

	outputPath := filepath.Join(configDir, "output.txt")
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("failed to save output.txt: %w", err)
	}

	configPath := filepath.Join(configDir, filepath.Base(cnf.CNFLoc))
	configJSON, err := json.MarshalIndent(cnf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to save config.json: %w", err)
	}

	// Database Insertion
	metadata := ArtifactMetadata{
		BuildID: cnf.Metadata.BuildId,
		Name:    cnf.Name(),

		// Version
		VersionMajor: cnf.Version().Major,
		VersionMinor: cnf.Version().Minor,
		VersionPatch: cnf.Version().Patch,

		// About Config
		BuildTime:          cnf.BuildDetails.BuildTime,
		BuildAt:            cnf.Metadata.BuildTime,
		Checksum:           cnf.ConfigHash(),
		ConfigFileLocation: configPath,
		Success:            cnf.BuildSuccess,

		Location: artifactDir,

		ArtifactName: cnf.BinaryObject.Name,
		ArtifactHash: cnf.ArtifactHash(),
	}

	result := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&metadata).Error; err != nil {
			return fmt.Errorf("Failed to save artifact metadata: %w", err)
		}

		return nil
	})

	if result != nil {
		return result
	}

	log.Printf("INFO :: Successfully saved artifact %s at %s\n", cnf.Name(), artifactDir)
	return nil
}

func (s *basicHandler) GetArtifactMetadata(buildId string) (*ArtifactMetadata, error) {
	var metadata ArtifactMetadata

	if err := s.db.Where("build_id = ?", buildId).First(&metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("artifact not found with build_id: %s", buildId)
		}
		return nil, fmt.Errorf("Failed to query artifact metadata: %w", err)
	}

	return &metadata, nil
}

func (s *basicHandler) GetArtifactConf(bId string) (*a2f.ArtifactConf, error) {
	metadata, err := s.GetArtifactMetadata(bId)
	if err != nil {
		return nil, err
	}

	// Read artifact binary
	artifactData, err := os.ReadFile(metadata.ConfigFileLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifact file: %w", err)
	}

	var a2f a2f.ArtifactConf
	if err := json.Unmarshal(artifactData, &a2f); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artifact config: %w", err)
	}

	return &a2f, nil
}

func (s *basicHandler) ListArtifacts() ([]BasicArtifactInfo, error) {
	var metadataList []ArtifactMetadata

	if err := s.db.Order("build_at DESC").Find(&metadataList).Error; err != nil {
		return nil, fmt.Errorf("failed to query artifacts: %w", err)
	}

	var artifacts []BasicArtifactInfo

	for _, metadata := range metadataList {

		atf := BasicArtifactInfo{
			BuildId:          metadata.BuildID,
			Name:             metadata.Name,
			Version:          a2f.Version{Major: metadata.VersionMajor, Minor: metadata.VersionMinor, Patch: metadata.VersionPatch},
			BuildTime:        metadata.BuildAt.Format(time.RFC3339),
			Success:          metadata.Success,
			ArtifactName:     metadata.ArtifactName,
			ConfigChecksum:   metadata.Checksum,
			ArtifactChecksum: metadata.ArtifactHash,
			BuildDuration:    metadata.BuildTime.String(),
		}

		artifacts = append(artifacts, atf)
	}

	return artifacts, nil
}
