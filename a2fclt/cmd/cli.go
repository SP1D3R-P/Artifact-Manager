package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	A2CLTCMD.CompletionOptions.DisableDefaultCmd = true

	A2CTL_BUILD_CMD.Flags().String("github", "", "GitHub repository URL to build the artifact from")

	A2CLTCMD.AddCommand(A2CTL_BUILD_CMD, A2CTL_INFO_CMD)
}

var A2CLTCMD = &cobra.Command{
	Use:   "a2clt",
	Short: "A2CLT is a command-line tool for automating the build ",
	Long: `
A2CLT is a command-line tool for automating the build artifact from source code.
It provides a simple and efficient way to build artifacts from source code, making it easier for developers to manage their build processes.
`,
}

var A2CTL_BUILD_CMD = &cobra.Command{
	Use:   "build [<source-path> | --github <github-repo>]",
	Short: "Build artifact from source code",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 1 {
			// Check if both source path and GitHub repository are provided
			if githubFlag, _ := cmd.Flags().GetString("github"); githubFlag != "" {
				fmt.Fprintf(os.Stderr, "ERROR :: Can't Mention both source path and github repository. Using %s for build.\n", args[0])
			}

			// Build from source path
			sourcePath := args[0]
			buildOptions := BuildOptions{
				SourcePath: sourcePath,
				BuildType:  BUILD_SC_FROM_PATH,
			}

			Build(buildOptions)

		} else if githubFlag, _ := cmd.Flags().GetString("github"); githubFlag != "" {

			// Build from GitHub repository [ IDEA :: it will fetch git clone shallow and then build the artifact ]
			buildOptions := BuildOptions{
				GithubRepo: githubFlag,
				BuildType:  BUILD_SC_FROM_GITHUB,
			}

			Build(buildOptions) // NOT IMPLEMENTED

		} else {
			fmt.Fprintln(os.Stderr, "ERROR :: Please provide either a source path or a GitHub repository for building the artifact.")
			cmd.Usage()
			os.Exit(1)
		}
	},
}

var A2CTL_INFO_CMD = &cobra.Command{
	Use:   "info",
	Short: "Show information about the build artifact",
	Long:  `Show information about the build artifact`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
