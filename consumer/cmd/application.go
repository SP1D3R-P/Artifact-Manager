package main

import (
	"context"
	"errors"
	"log"
	"path/filepath"

	a2f_manager "example.com/consumer/internal/manager"
	Watcher "example.com/consumer/internal/watchers"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
	watchPaths []string
}

type Application struct {
	cnf       *Config
	manager   *a2f_manager.Manager
	fsWatcher Watcher.FSWatcher
}

func NewApplication(cnf *Config, manager *a2f_manager.Manager) Application {

	watcher, err := Watcher.NewFSWatcher()

	if err != nil {
		log.Fatalf("Couldn't Start Application Due to : %s", err.Error())
	}

	watcher.AddPaths(cnf.watchPaths...)

	artifactManagerWorkflow := newWorkflowCallbackFn(manager)
	watcher.AddCallbacks(fsnotify.Create, artifactManagerWorkflow)

	return Application{
		cnf:       cnf,
		manager:   manager,
		fsWatcher: watcher,
	}
}

func (app *Application) Run(ctx context.Context) {
	defer app.fsWatcher.Close()

	// blocking work so must call with corutine
	app.fsWatcher.Watch(ctx)

}

///////////////////////////////////////////////////////////////////////////////
//							 CallBack Fn
///////////////////////////////////////////////////////////////////////////////

type workflowCallbabckFn struct {
	manager *a2f_manager.Manager
}

func newWorkflowCallbackFn(manager *a2f_manager.Manager) *workflowCallbabckFn {
	return &workflowCallbabckFn{
		manager: manager,
	}
}

func (wf *workflowCallbabckFn) Call(arg interface{}) error {
	file, ok := arg.(string)
	if !ok {
		return errors.New("Invalid arg is passed")
	}

	log.Printf("New File Created %s\n", file)

	// Not looking for Other format
	if filepath.Ext(file) != ".json" {
		return nil
	}

	artifact, err := wf.manager.LoadArtifact(file)
	if err != nil {
		return err
	}

	if err := wf.manager.ProcessArtifact(artifact); err != nil {
		return err
	}

	err = wf.manager.SaveArtifact(artifact)

	return err
}
