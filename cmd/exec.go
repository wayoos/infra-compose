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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wayoos/infra-compose/compose"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a global or service command.",
	Long:  `Run a global or service command.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("exec args: " + strings.Join(args, " "))

		err := cmd.Flags().Parse(args)
		if err != nil {
			fmt.Println(err)
		}

		if args[0] == "--help" || args[0] == "-h" {
			cmd.HelpFunc()(cmd, args)
			return nil
		}

		//			cmd.
		fmt.Println("compose file: " + composeFile)
		fmt.Println("compose file: " + projectDir)

		fmt.Println("cmd args: " + strings.Join(cmd.Flags().Args(), " "))

		compose := compose.Compose{}
		err = compose.Load(composeFile, projectDir)
		if err != nil {
			return err
		}

		return compose.Exec(args)
	},
	DisableFlagParsing: true,
	SilenceUsage:       true,
	SilenceErrors:      true,
}

func init() {
	RootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
