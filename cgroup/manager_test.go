package cgroup

import "testing"

func TestNewManager(t *testing.T) {
	m, err := NewManager("/sys/fs/cgroup", "/axliu", &Resources{})
	if err != nil {
		t.Errorf("new manager error: %v", err)
		return
	}
	t.Logf("manager unifiedMountpoint: %s", m.unifiedMountpoint)
	t.Logf("manager path: %s", m.path)
}
