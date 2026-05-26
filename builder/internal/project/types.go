package project

type Step struct {
	Cmd   string   `json:"cmd"`
	Input []string `json:"input,omitempty"`
}

type ProcessInfo struct {
	Steps []Step            `json:"steps"`
	Env   map[string]string `json:"environ,omitempty"`
}

type projectConfig struct {

	// project name
	Project string `json:"project"`

	// version of the project
	Version string `json:"version"`

	// where the artifact will be stored
	StorageLocation string `json:"location"`

	// name of the artifact [ finally be stored ]
	Artifact string `json:"artifact"`

	/*
		How the Project Behave [ building and execution ]
	*/

	// building process
	Build ProcessInfo `json:"build"`
	// execution process
	Exec ProcessInfo `json:"exec"`
}

type Project struct {
	location string
	config   *projectConfig

	// BUILD_VERSION && PROJECT_NAME
	envs map[string]string

	// This is the unique id for the build process
	BuildId string
}

type Version struct {
	Major int
	Minor int
	Patch int
}
