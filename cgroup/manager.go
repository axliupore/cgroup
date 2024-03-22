package cgroup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	unifiedMountpoint string
	path              string
}

// NewManager new a cgroup, and set resources controllers
func NewManager(mountpoint string, group string, resources *Resources) (*Manager, error) {
	if resources == nil {
		return nil, errors.New("resources reference is nil")
	}
	if err := VerifyGroupPath(group); err != nil {
		return nil, err
	}
	path := filepath.Join(mountpoint, group)
	if err := os.MkdirAll(path, dirPerm); err != nil {
		return nil, err
	}
	m := &Manager{
		unifiedMountpoint: mountpoint,
		path:              path,
	}

	if err := m.ToggleControllers(resources.EnabledControllers(), Enable); err != nil {
		// clean up cgroup dir on failure
		os.Remove(path)
		return nil, err
	}
	if err := setResources(path, resources); err != nil {
		os.Remove(path)
		return nil, err
	}
	return m, nil
}

func (m *Manager) ToggleControllers(controllers []string, t ControllerToggle) error {
	// when m.path is like /foo/bar/baz, the following files need to be written:
	// * /sys/fs/cgroup/cgroup.subtree_control
	// * /sys/fs/cgroup/foo/cgroup.subtree_control
	// * /sys/fs/cgroup/foo/bar/cgroup.subtree_control
	// Note that /sys/fs/cgroup/foo/bar/baz/cgroup.subtree_control does not need to be written.
	split := strings.Split(m.path, "/")
	var lastErr error
	for i := range split {
		f := strings.Join(split[:i], "/")
		if !strings.HasPrefix(f, m.unifiedMountpoint) || f == m.path {
			continue
		}
		filePath := filepath.Join(f, subtreeControl)
		if err := m.writeSubtreeControl(filePath, controllers, t); err != nil {
			lastErr = fmt.Errorf("failed to write subtree controllers %+v to %q: %w", controllers, filePath, err)
		} else {
			lastErr = nil
		}
	}
	return lastErr
}

func (m *Manager) writeSubtreeControl(path string, controllers []string, t ControllerToggle) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		fmt.Printf("err: %v", err)
		return err
	}
	defer f.Close()
	if t == Enable {
		controllers = toggleFunc(controllers, "+")
	} else {
		controllers = toggleFunc(controllers, "-")
	}
	_, err = f.WriteString(strings.Join(controllers, " "))
	return err
}

func toggleFunc(controllers []string, prefix string) []string {
	out := make([]string, len(controllers))
	for i, v := range controllers {
		out[i] = prefix + v
	}
	return out
}

// RootControllers read /sys/fs/cgroup/cgroup.controllers content
// [cpuset cpu io memory pids rdma hugetlb]
func (m *Manager) RootControllers() ([]string, error) {
	b, err := os.ReadFile(filepath.Join(m.unifiedMountpoint, controllersFile))
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(b)), nil
}
