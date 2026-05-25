package main

import (
	"context"
	"log"
	"os"
	"sync"

	"example.com/consumer/cmd/routes"
	"example.com/consumer/internal"
	a2f_manager "example.com/consumer/internal/manager"
	"example.com/consumer/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	if _, err := os.Stat(internal.StorageFSDir); err != nil {
		log.Fatalf("Could Not Get Storage Dir due to Error : %s\n", err)
	}

	if _, err := os.Stat(internal.ArtifactDir); err != nil {
		log.Fatalf("Could Not Get artifact Dir due to Error : %s\n", err)
	}

	log.Printf("Path to Artifacts : %s\n", internal.ArtifactDir)
	log.Printf("Path to Storage : %s\n", internal.StorageFSDir)

	cnf := Config{
		watchPaths: []string{
			internal.ArtifactDir,
		},
	}
	storageHandler := storage.NewBasicHandler()
	m := a2f_manager.NewManager(storageHandler)
	app := NewApplication(&cnf, m)

	wg := sync.WaitGroup{}

	ctx := context.Background()
	// ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	// Starting file watcher
	wg.Go(func() {
		app.Run(ctx)
	})

	// Start API server in goroutine

	wg.Go(func() {
		startAPIServer(m)
	})

	wg.Wait()
}

func startAPIServer(manager *a2f_manager.Manager) {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup routes
	apiHandler := routes.NewAPIHandler(manager, internal.StorageFSDir)
	apiHandler.SetupRoutes(router)

	// Start server on port 8080
	log.Println("Starting API server on :9000")
	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Failed to start API server: %v\n", err)
	}
}
