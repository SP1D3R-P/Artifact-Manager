package artifact

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ArtifactWriter interface {
	SaveTo(loc string, a2f *Artifact) error
}

type artifactWriter struct {
}

func NewBaseArtifactWriter() ArtifactWriter {
	return &artifactWriter{}
}

func (w *artifactWriter) SaveTo(dir string, a2f *Artifact) error {

	dumpingDir := filepath.Join(dir, a2f.BuildDetails.BuildOutput)
	fmt.Printf("INFO :: Dumping the artifact to the location: %s\n", dumpingDir)

	if err := os.MkdirAll(dumpingDir, 0755); err != nil {
		return err
	}

	// coping the binary
	{
		// destination path
		DestBinPath := filepath.Join(dumpingDir, a2f.BinaryObject.Name)
		// a2f in the project generated
		SourceBinPath := filepath.Join(a2f.Metadata.sourcePath, a2f.BinaryObject.Name)

		bytesRead, err := os.ReadFile(SourceBinPath)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("The Expected Artifact[%s] is not Generated During the process", a2f.BinaryObject.Name)
			}
			return nil
		}

		if err := os.WriteFile(DestBinPath, bytesRead, 0755); err != nil {
			return err
		}
	}

	// writing the ouput output [ for future refrence ig ]
	{
		OuputDest := filepath.Join(dumpingDir, "output.txt")
		if err := os.WriteFile(OuputDest, []byte(a2f.result.BuildOutput), 0644); err != nil {
			return err
		}
	}

	// writing the config to the [dir] => artifact/
	// this will be the first thing that the consumer will load
	configPath := filepath.Join(dir, a2f.Filename())
	if err := os.WriteFile(configPath, a2f.Data(), 0644); err != nil {
		return err
	}

	return nil
}
