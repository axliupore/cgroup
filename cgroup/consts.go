package cgroup

const (
	unifiedMountpoint = "/sys/fs/cgroup" // unified mount point in file system

	subtreeControl = "cgroup.subtree_control"
	controllersFile = "cgroup.controllers"
	cgroupProcs    = "cgroup.procs"

	initPath = "init"

	filePerm = 0644
	dirPerm = 0755
)

type CGMode int

const (
	Unified     CGMode = iota // Unified with only cgroup v2 mounted
	Legacy                    // Legacy cgroup v1
	Hybrid                    // Hybrid with cgroup v1 and v2 controllers mounted
	Unavailable               // Unavailable cgroup mountpoint
)

type ControllerToggle int

const (
	Enable ControllerToggle = iota + 1
	Disable
)
