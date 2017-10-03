// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wayoos/infra-compose/compose"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [flags] [ENVIRONMENT] [SERVICE] COMMAND [ARGS...]",
	Short: "Run a global or service command",
	Long:  `Run a global or service command`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExecute(cmd, args, false)
	},
	DisableFlagParsing: true,
	SilenceUsage:       true,
	SilenceErrors:      true,
}

func runExecute(cmd *cobra.Command, args []string, dryRun bool) error {
	// Parse flag and ignore invalid flag used in sub command
	cmd.Flags().Parse(args)

	if args[0] == "--help" || args[0] == "-h" {
		cmd.HelpFunc()(cmd, args)
		return nil
	}

	compose := compose.Compose{}
	compose.DryRun = dryRun
	err := compose.Load(composeFile, projectDir)
	if err != nil {
		return err
	}

	// remove global flags
	firstCmdArg := cmd.Flags().Arg(0)
	validArgs := args
	for firstCmdArg != validArgs[0] {
		validArgs = validArgs[1:]
	}

	return compose.Exec(validArgs)
}

func init() {
	RootCmd.AddCommand(execCmd)
}
