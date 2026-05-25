// pronounciation as build for generic
package build4generic

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"example.com/builder/internal"
	"example.com/builder/internal/builder"
	"example.com/builder/internal/project"
)

type genericBuilder struct {
	Project *project.Project
}

func NewBuilder(proj *project.Project) builder.Builder {
	return &genericBuilder{Project: proj}
}

func (b *genericBuilder) Build() *builder.BuildResult {
	result := &builder.BuildResult{
		ProjectName: b.Project.Name(),
	}

	// building details
	builder := b.Project.BuildProcess()

	// request for envs
	envs := internal.ResolveENV(internal.EnvSliceToMap(os.Environ()), builder.Env, b.Project.BaseEnvs())

	startTime := time.Now()

	Ouputs := make([]string, 0)
	defer func() {
		result.BuildOutput += fmt.Sprintf(":{%%CONSOLE%%}:\n%s\n", strings.Join(Ouputs, "\n"))
	}()

	for _, step := range builder.Steps {

		// replace the string quote to normal string
		cmd := strings.ReplaceAll(step.Cmd, `/"`, `"`)
		log.Println("CMD :: ", cmd)

		var in = []string{} // transformed inputs
		for _, input := range step.Input {
			in = append(
				in,
				os.Expand(strings.ReplaceAll(input, `/"`, `"`), envs),
			)
		}

		// expand the enviroment variables
		cmd = os.Expand(cmd, envs)

		// execute and append the output
		result.BuildOutput += fmt.Sprintf(":{%%CMD%%}:  %s\n", cmd)
		output, err := b.executeCommand(cmd, in...)
		if output != "" {
			Ouputs = append(Ouputs, output)
		}

		if err != nil {
			result.BuildTime = time.Duration(time.Since(startTime).Milliseconds())
			result.Error = fmt.Errorf("build failed at step '%s': %w", step, err)
			result.Success = false
			return result
		}
	}

	result.BuildTime = time.Duration(time.Since(startTime).Milliseconds())

	// checking if the artifact is generated after the build process or not
	artifactPath := b.Project.ArtifactLocation()
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("artifact not found at %s", artifactPath)
		result.Success = false
		return result
	}

	result.ArtifactPath = artifactPath
	result.Success = true
	return result
}

func (b *genericBuilder) executeCommand(cmdStr string, in ...string) (string, error) {

	var cmd *exec.Cmd
	var newCmdStr []string

	if isWindowsCommand() {
		newCmdStr = append(newCmdStr, "/c", cmdStr)
		newCmdStr = append(newCmdStr, in...)
		cmd = exec.Command("cmd", newCmdStr...)
	} else {
		newCmdStr = append(newCmdStr, "-c", cmdStr)
		newCmdStr = append(newCmdStr, in...)
		cmd = exec.Command("sh", newCmdStr...)
	}

	// adding the env
	cmd.Env = append(os.Environ(), internal.MapToEnvSlice(b.Project.BuildProcess().Env)...)
	cmd.Env = append(cmd.Env, internal.MapToEnvSlice(b.Project.BaseEnvs())...)

	// working dir
	cmd.Dir = b.Project.ProjectLocation()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()

	if stderr.Len() > 0 {
		output += "\nSTDERR: " + stderr.String()
	}

	return output, err
}

func isWindowsCommand() bool {
	return os.PathSeparator == '\\'
}
