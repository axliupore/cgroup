package cgroup

import (
	"strconv"
	"strings"
)

// CPUMax represents the maximum CPU quota in a cgroup.
// It is a string type alias used to define the maximum CPU quota, which can be either a numeric value (representing microseconds) or "max" (indicating no limit).
type CPUMax string

// CPU represents the CPU-related attributes in a cgroup.
type CPU struct {
	Weight *uint64 // Weight is a pointer to a uint64, indicating the CPU weight for relative priority in resource allocation.
	Max    CPUMax  // Max represents the maximum CPU quota in the cgroup.
	Cpus   string  // Cpus represents the list of available CPU cores in the CPU group.
	Mems   string  // Mems represents the list of available memory nodes in the CPU group.
}

func NewCPUMax(quota *int64, period *uint64) CPUMax {
	max := "max"
	if quota != nil {
		max = strconv.FormatInt(*quota, 10)
	}
	return CPUMax(strings.Join([]string{max, strconv.FormatUint(*period, 10)}, " "))
}

func (c *CPU) Values() (v []Value) {
	if c.Weight != nil {
		v = append(v, Value{
			filename: "cpu.weight",
			value:    *c.Weight,
		})
	}
	if c.Max != "" {
		v = append(v, Value{
			filename: "cpu.max",
			value:    c.Max,
		})
	}
	if c.Cpus != "" {
		v = append(v, Value{
			filename: "cpuset.cpus",
			value:    c.Cpus,
		})
	}
	if c.Mems != "" {
		v = append(v, Value{
			filename: "cpuset.mems",
			value:    c.Mems,
		})
	}
	return
}
