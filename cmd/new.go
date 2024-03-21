package main

import (
	"fmt"
	"github.com/axliupore/cgroup/cgroup"
	"github.com/spf13/cobra"
	"strconv"
)

var newCommand = &cobra.Command{
	Use:   "new path",
	Short: "create a new cgroup",
	Long:  "if multiple parameters are entered, only the first one will be created by default",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("input new cgroup path")
			return
		}
		// if multiple parameters are entered, only the first one will be created by default
		path := args[0]
		mountpoint := cmd.Parent().PersistentFlags().Lookup("mountpoint").Value.String()
		m, err := cgroup.NewManager(mountpoint, path, &cgroup.Resources{})
		if err != nil {
			fmt.Println(err)
		}
		enable, _ := strconv.ParseBool(cmd.Flags().Lookup("enable").Value.String())
		if enable {
			controllers, err := m.RootControllers()
			if err != nil {
				fmt.Println(err)
			}
			if err := m.ToggleControllers(controllers, cgroup.Enable); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	newCommand.Flags().BoolP("enable", "e", false, "enable the controllers for the group")
}
