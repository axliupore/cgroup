package cgroup

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestCgroupCPUStats(t *testing.T) {
	group := "/cpu-test"
	groupPath := fmt.Sprintf("%s-%d", group, os.Getpid())
	var (
		qota   int64  = 10000
		period uint64 = 8000
		weight uint64 = 100
	)
	max := "10000 8000"
	res := &Resources{
		CPU: &CPU{
			Weight: &weight,
			Max:    NewCPUMax(&qota, &period),
			Cpus:   "0",
			Mems:   "0",
		},
	}
	m, err := NewManager(unifiedMountpoint, groupPath, res)
	require.NoError(t, err, "failed to init new cgroup manager")
	t.Cleanup(func() {
		os.Remove(m.path)
	})
	checkFileContent(t, m.path, "cpu.weight", strconv.FormatUint(weight, 10))
	checkFileContent(t, m.path, "cpu.max", max)
	checkFileContent(t, m.path, "cpuset.cpus", "0")
	checkFileContent(t, m.path, "cpuset.mems", "0")
}
