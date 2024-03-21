package cgroup

import "testing"

import "fmt"

func TestMode(t *testing.T) {
	mode := Mode()
	if mode == Unified {
		fmt.Println("cgroup: unified [cgroup mode is cgroup v2]")
	} else if mode == Legacy {
		fmt.Println("cgroup: legacy [cgroup mode is cgroup v1]")
	} else if mode == Hybrid {
		fmt.Println("cgroup: hybrid [cgroup mode is mixed cgroup v1 and v2]")
	} else {
		fmt.Println("cgroup: unavailable [cgroup mode detection failed]")
	}
}

func TestEnableNesting(t *testing.T) {
	if err := EnableNesting(); err != nil {
		fmt.Printf("err: %v", err)
	}
}
