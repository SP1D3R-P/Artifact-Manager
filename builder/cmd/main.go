package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"example.com/builder/internal/artifact"
	build4generic "example.com/builder/internal/builder/generic"
	"example.com/builder/internal/project"
)

const (
	artifactsDir = "/artifacts"
	codeDir      = "/code"
	MiB          = 1024 * 1024
	MaxFileSize  = 50 * MiB // 50 MB
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/build/compressed", HandleBuildRequest)

	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check called")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

type BuildResponse struct {
	BuildID string `json:"build_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

func HandleBuildRequest(w http.ResponseWriter, r *http.Request) {

	writeResponse := func(statusCode int, response BuildResponse) {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}

	limitedReader := io.LimitReader(r.Body, MaxFileSize)
	defer r.Body.Close()

	gzReader, err := gzip.NewReader(limitedReader)
	if err != nil {
		writeResponse(http.StatusBadRequest, BuildResponse{Error: "Invalid gzip format"})
		return
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	randID := make([]byte, 16)
	rand.Read(randID)
	BuildId := fmt.Sprintf("Build-%x", randID)

	tempDir := filepath.Join(codeDir, BuildId)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		writeResponse(http.StatusInternalServerError, BuildResponse{Error: "Failed to create temp directory"})
		return
	}

	// Unpacking the tar.gz data to tempDir
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			writeResponse(http.StatusInternalServerError, BuildResponse{Error: fmt.Sprintf("Failed to read tar: %v", err)})
			return
		}

		target := filepath.Join(tempDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				continue
			}
			io.Copy(f, tarReader)
			f.Close()
		}
	}

	go func() {
		if err := run(tempDir, BuildId); err != nil {
			log.Printf("ERROR :: Build failed: %v\n", err)
		}
		// Clean up the temporary directory after the build is done
		os.RemoveAll(tempDir)
	}()

	writeResponse(http.StatusOK, BuildResponse{BuildID: BuildId})

}

func run(wd, buildId string) error {

	// validating the directory and if config file present
	err := os.Chdir(wd)
	if err != nil {
		return fmt.Errorf("Failed to switch working dir")
	}

	cnf := filepath.Join(wd, "config.json")
	if _, err := os.Stat(cnf); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Config File Doesn't Exists.")
		}
		return err
	}

	// building the project
	if err := buildProject(wd, buildId); err != nil {
		return fmt.Errorf("Failed to build project %s: %v\n", wd, err)
	}
	return nil
}

func buildProject(projPath, buildId string) error {

	log.Printf("INFO :: Starting building [Build ID: %s]\n", buildId)
	proj, err := project.LoadProject(projPath, buildId)
	if err != nil {
		return fmt.Errorf("Failed to load project: %w", err)
	}

	log.Printf("INFO :: Building project: %s (v%v)\n", proj.Name(), proj.Version())

	bldr := build4generic.NewBuilder(proj)
	result := bldr.Build()

	if !result.Success {
		log.Printf("ERROR :: Build failed: %v\n", result.Error)
		return result.Error
	}

	log.Println("INFO :: Build successful.")
	log.Printf("BUILD OUTPUT :: \n%s\n", result.BuildOutput)

	art, err := artifact.NewArtifact(proj, result)
	if err != nil {
		return fmt.Errorf("failed to create artifact: %w\n", err)
	}

	if err := art.SaveTo(artifactsDir); err != nil {
		return fmt.Errorf("failed to save artifact: %w\n", err)
	}

	log.Printf("Artifact saved to: %s/%s\n", artifactsDir, art.Filename())
	if art.BuildDetails.BuildOutput != "" {
		log.Printf("Artifact also copied to location: %s\n", art.BuildDetails.BuildOutput)
	}

	return nil
}
