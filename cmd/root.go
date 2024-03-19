package main

import "github.com/spf13/cobra"

// execute this program
func execute() {
	rootCmd := &cobra.Command{
		Use:   "cgruop",
		Short: "cgroup managing cgroups tool",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	// add commands
	rootCmd.AddCommand(mode)
	rootCmd.Execute()
}
