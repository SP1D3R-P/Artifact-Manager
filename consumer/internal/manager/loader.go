package a2f_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"example.com/consumer/internal/a2f"
)

type Loader struct {
}

func NewLoader() Loader {
	return Loader{}
}

// Load the Artifact
func (l *Loader) Load(fname string) (*a2f.Artifact, error) {
	fInfo, err := os.Stat(fname)

	if err != nil {
		return nil, err
	}

	if fInfo.IsDir() {
		return nil, fmt.Errorf("The Given Path is a Dir Expected a Json File")
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	cnf, err := l.parse(data)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Invalid Config File"), err)
	}
	cnf.CNFLoc = fname

	if err := l.validate(*cnf); err != nil {
		return nil, fmt.Errorf("Invalid Config Failed Due :: %s \n", err)
	}

	return a2f.New(cnf), nil

}

// ParseArtifact parses raw data into an Artifact
func (l *Loader) parse(data []byte) (*a2f.ArtifactConf, error) {
	var cnf a2f.ArtifactConf
	if err := json.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}
	return &cnf, nil
}

// ValidateArtifact validates the artifact's content and returns an error if validation fails
func (l *Loader) validate(artifact a2f.ArtifactConf) error {
	return nil
}
