package main

import (
	"fmt"
	"github.com/axliupore/cgroup/cgroup"
	"github.com/spf13/cobra"
)

var mode = &cobra.Command{
	Use:   "mode",
	Short: "return the cgroup mode that is mounted on the system",
	Run: func(cmd *cobra.Command, args []string) {
		mode := cgroup.Mode()
		if mode == cgroup.Unified {
			fmt.Println("cgroup: unified [cgroup mode is cgroup v2]")
		} else if mode == cgroup.Legacy {
			fmt.Println("cgroup: legacy [cgroup mode is cgroup v1]")
		} else if mode == cgroup.Hybrid {
			fmt.Println("cgroup: hybrid [cgroup mode is mixed cgroup v1 and v2]")
		} else {
			fmt.Println("cgroup: unavailable [cgroup mode detection failed]")
		}
	},
}
