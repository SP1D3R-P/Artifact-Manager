package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	A2CTL_BASE_URL = "http://localhost:8080/api/v1"
)

func main() {
	A2CLTCMD.Execute()

}

const (
	BUILD_SC_FROM_PATH = iota
	BUILD_SC_FROM_GITHUB
)

type BuildOptions struct {
	SourcePath string
	GithubRepo string
	BuildType  int
}

type StatusCode int
type BuildResponse struct {
	BuildID string `json:"build_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Build(options BuildOptions) {

	switch options.BuildType {
	case BUILD_SC_FROM_PATH:
		{
			fmt.Printf("INFO :: Building artifact from source path: %s\n", options.SourcePath)

			// Validate the source path
			fileStat, err := os.Stat(options.SourcePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR :: Error accessing source path: %v\n", err)
				os.Exit(1)
			}
			if !fileStat.IsDir() {
				fmt.Fprintf(os.Stderr, "ERROR :: Provided source path is not a directory: %s\n", options.SourcePath)
				os.Exit(1)
			}

			// Compress the source code into a tar.gz format
			compressedSCData := bytes.NewBuffer(nil)
			compressedSC := compressFolder(options.SourcePath, compressedSCData)

			if compressedSC != nil {
				fmt.Fprintf(os.Stderr, "ERROR :: Failed to compress source code due : %v\n", compressedSC)
				os.Exit(1)
			}

			// Upload the compressed source code to the server
			statCode, resp := uploadCompressedSC(compressedSCData.Bytes())
			if statCode != http.StatusOK {
				fmt.Fprintf(os.Stderr, "ERROR :: Failed to upload compressed source code: %v\n", resp.Error)
				os.Exit(1)
			}
			fmt.Printf("INFO :: Upload successful.\n")
			fmt.Printf("INFO :: The Build ID is %s\n", resp.BuildID)
		}

	case BUILD_SC_FROM_GITHUB:
		{
			fmt.Printf("INFO :: Building artifact from GitHub repository: %s\n", options.GithubRepo)
			panic("Building from GitHub repository is not implemented yet.")
		}
	default:
		// Unreachable code
		fmt.Fprintf(os.Stderr, "Invalid build type specified.\n")
		os.Exit(1)
	}
}

// Compress Folder to tar.gz format
func compressFolder(src string, buf io.Writer) error {

	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, _ := tar.FileInfoHeader(info, "")
		header.Name, _ = filepath.Rel(filepath.Dir(src), path)
		tw.WriteHeader(header)

		if !info.IsDir() {
			f, _ := os.Open(path)
			defer f.Close()
			io.Copy(tw, f)
		}
		return nil
	})
}

func uploadCompressedSC(data []byte) (StatusCode, BuildResponse) {
	var resp *http.Response

	// Upload the compressed source code to the server
	BuildLink := fmt.Sprintf("%s/build/compressed", A2CTL_BASE_URL)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", BuildLink, bytes.NewBuffer(data))
	resp, err = client.Do(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :: Failed to create HTTP request: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	rsp_byt_data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :: Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	var rsp BuildResponse

	if err := json.Unmarshal(rsp_byt_data, &rsp); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :: Failed to unmarshal response: %v\n", err)
		fmt.Fprintf(os.Stderr, "Response body: %s\n", string(rsp_byt_data))
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		return StatusCode(resp.StatusCode), rsp
	}
	return StatusCode(resp.StatusCode), rsp
}
