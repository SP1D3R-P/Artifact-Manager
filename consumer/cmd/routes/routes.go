package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	a2f_manager "example.com/consumer/internal/manager"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	manager     *a2f_manager.Manager
	storagePath string
}

func NewAPIHandler(manager *a2f_manager.Manager, storagePath string) *APIHandler {
	return &APIHandler{
		manager:     manager,
		storagePath: storagePath,
	}
}

func (h *APIHandler) SetupRoutes(router *gin.Engine) {

	router.GET("/builds", h.listArtifacts)
	router.GET("/builds/:BuildId", h.getArtifactConf)

	router.GET("/:BuildId/output.txt", h.getOutputFile)

	router.GET("/artifacts/*artifact_path", h.getArtifactFile)
}

// listArtifacts returns all stored artifacts
func (h *APIHandler) listArtifacts(c *gin.Context) {
	artifacts, err := h.manager.ListArtifacts()
	if err != nil {
		log.Printf("Error listing artifacts: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list artifacts",
		})
		return
	}

	if len(artifacts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"artifacts": []interface{}{},
			"count":     0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"artifacts": artifacts,
		"count":     len(artifacts),
	})
}

func (h *APIHandler) getArtifactConf(c *gin.Context) {
	buildNo := c.Param("BuildId")
	if buildNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BuildId is required",
		})
		return
	}
	artifactConf, err := h.manager.GetArtifactConf(buildNo)
	if err != nil {
		log.Printf("ERROR :: Error getting artifact config %s: %v\n", buildNo, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Artifact config not found: %s", buildNo),
		})
		return
	}

	c.JSON(http.StatusOK, artifactConf)
}

// getOutputFile returns the output.txt file for a specific artifact - BuildId
func (h *APIHandler) getOutputFile(c *gin.Context) {
	buildNo := c.Param("BuildId")

	if buildNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BuildId is required",
		})
		return
	}

	// Try to find artifact by name first
	output, err := h.manager.GetArtifactOutput(buildNo)
	if err != nil {
		log.Printf("ERROR :: Error getting artifact output %s: %v\n", buildNo, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Artifact output not found: %s", buildNo),
		})
		return
	}

	c.Data(http.StatusOK, "text/plain", output)
}

// getArtifactFile returns the artifact file (binary) for a specific artifact
func (h *APIHandler) getArtifactFile(c *gin.Context) {
	artifactPath := c.Param("artifact_path")
	if artifactPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Artifact path is required",
		})
		return
	}

	if filepath.Base(artifactPath) == "output.txt" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Use /:BuildId/output.txt to access output.txt files",
		})
		return
	}

	artifactFilePath := filepath.Join(h.storagePath, "artifacts", artifactPath)
	if _, err := os.Stat(artifactFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Artifact file not found",
		})
		return
	}
	c.File(artifactFilePath)
}
