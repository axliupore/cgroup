package cgroup

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

var (
	modeOnce sync.Once // make sure only check cgroup mode once
	cgMode   CGMode    // storage system's cgroup mode
)

// Mode returns the cgroup mode running the host.
func Mode() CGMode {
	modeOnce.Do(func() {
		var st unix.Statfs_t
		// obtain file system statistics for a specified path
		if err := unix.Statfs(unifiedMountpoint, &st); err != nil {
			cgMode = Unavailable
			return
		}
		// support and use only cgroup v2
		if st.Type == unix.CGROUP2_SUPER_MAGIC {
			cgMode = Unified
		} else {
			cgMode = Legacy
			// only cgroup v1
			if err := unix.Statfs(filepath.Join(unifiedMountpoint, "unified"), &st); err != nil {
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

func readFile(p string) ([]byte, error) {
	data, err := os.ReadFile(p)
	if err != nil && errors.Is(err, syscall.EINTR) {
		data, err = os.ReadFile(p)
	}
	return data, err
}

// EnableV2Nesting migrates all process in the container to nested /init path
// and enables all available controllers in the root cgroup
func EnableNesting() error {
	if Mode() != Unified {
		return nil
	}

	p, err := readFile(filepath.Join(unifiedMountpoint, cgroupProcs))
	if err != nil {
		return err
	}
	// read all process from cgroup.procs
	procs := strings.Split(string(p), "\n")
	if len(procs) == 0 {
		return nil
	}

	// mkdir init
	if err := os.Mkdir(filepath.Join(unifiedMountpoint, initPath), dirPerm); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	// move all process into init cgroup
	procFile, err := os.OpenFile(filepath.Join(unifiedMountpoint, initPath, cgroupProcs), os.O_RDWR, filePerm)
	if err != nil {
		return err
	}
	for _, v := range procs {
		if _, err := procFile.WriteString(v); err != nil {
			continue
		}
	}
	procFile.Close()
	return nil
}
