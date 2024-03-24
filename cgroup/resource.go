package cgroup

import "github.com/opencontainers/runtime-spec/specs-go"

// Resources for a cgruop unified hierarchy
type Resources struct {
	CPU     *CPU
	Memory  *Memory
	Devices []specs.LinuxDeviceCgroup // When len(Devices) is zero, devices are not controlled
}

// EnabledControllers get resources enabled controllers
func (r *Resources) EnabledControllers() (c []string) {
	if r.CPU != nil {
		c = append(c, "cpu")
		if r.CPU.Cpus != "" || r.CPU.Mems != "" {
			c = append(c, "cpuset")
		}
	}
	if r.Memory != nil {
		c = append(c, "memory")
	}
	return
}

func (r *Resources) Values() (v []Value) {
	if r.CPU != nil {
		v = append(v, r.CPU.Values()...)
	}
	if r.Memory != nil {
		v = append(v, r.Memory.Values()...)
	}
	return
}

func setResources(path string, resources *Resources) error {
	if resources != nil {
		if err := writeValues(path, resources.Values()); err != nil {
			return err
		}
		if err := setDevices(path, resources.Devices); err != nil {
			return err
		}
	}
	return nil
}
