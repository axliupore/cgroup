package cgroup

import (
	"golang.org/x/sys/unix"
	"path/filepath"
	"sync"
)

var (
	modeOnce sync.Once // make sure only check cgroup mode once
	cgMode   CGMode    // storage system's cgroup mode
)

// unified mount point in file system
const unifiedMountPoint = "/sys/fs/cgroup"

type CGMode int

const (
	Unified     CGMode = iota // Unified with only cgroup v2 mounted
	Legacy                    // Legacy cgroup v1
	Hybrid                    // Hybrid with cgroup v1 and v2 controllers mounted
	Unavailable               // Unavailable cgroup mountpoint
)

// Mode returns the cgroup mode running the host.
func Mode() CGMode {
	modeOnce.Do(func() {
		var st unix.Statfs_t
		// obtain file system statistics for a specified path
		if err := unix.Statfs(unifiedMountPoint, &st); err != nil {
			cgMode = Unavailable
			return
		}
		// support and use only cgroup v2
		if st.Type == unix.CGROUP2_SUPER_MAGIC {
			cgMode = Unified
		} else {
			cgMode = Legacy
			// only cgroup v1
			if err := unix.Statfs(filepath.Join(unifiedMountPoint, "unified"), &st); err != nil {
				return
			}
			// support both cgroup v1 and cgroup v2
			if st.Type == unix.CGROUP2_SUPER_MAGIC {
				cgMode = Hybrid
			}
		}
	})
	return cgMode
}
