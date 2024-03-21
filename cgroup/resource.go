package cgroup

// Resources for a cgruop unified hierarchy
type Resources struct {
	CPU *CPU
}

func (r *Resources) EnabledControllers() (c []string) {
	if r.CPU != nil {
		c = append(c, "cpu")
		if r.CPU.Cpus != "" || r.CPU.Mems != "" {
			c = append(c, "cpuset")
		}
	}
	return
}

func (r *Resources) Values() (v []Value) {
	if r.CPU != nil {
		v = append(v, r.CPU.Values()...)
	}
	return
}

func setResources(path string, resources *Resources) error {
	if resources != nil {
		if err := writeValues(path, resources.Values()); err != nil {
			return err
		}
	}
	return nil
}