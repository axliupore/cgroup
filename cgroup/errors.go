package cgroup

import "errors"

var (
	ErrInvalidGroupPath = errors.New("cgroups: invalid group path")
	ErrInvalidFormat    = errors.New("cgroups: parsing file with invalid format failed")
)
