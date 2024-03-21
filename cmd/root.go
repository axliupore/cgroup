package main

import (
	"github.com/axliupore/cgroup/cgroup"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "cgruop",
	Short: "cgroup managing cgroups tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// execute this program
func execute() {
	// add commands
	rootCommand.AddCommand(modeCommand, newCommand)
	rootCommand.Execute()
}

func init() {
	err := cgroup.EnableNesting()
	if err != nil {
		panic(err)
	}
	rootCommand.PersistentFlags().StringP("mountpoint", "m", "/sys/fs/cgroup", "cgroup mountpoint")
}
