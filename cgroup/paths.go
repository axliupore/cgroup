package cgroup

import (
	"path/filepath"
	"strings"
)

// VerifyGroupPath verifies the format of group path string p.
// The format is same as the third field in /proc/PID/cgroup.
// e.g. "/user.slice/user-1001.slice/session-1.scope"
//
// p must be a "clean" absolute path starts with "/", and must not contain "/sys/fs/cgroup" prefix.
//
// VerifyGroupPath doesn't verify whether p actually exists on the system.
func VerifyGroupPath(p string) error {
	if !strings.HasPrefix(p, "/") {
		return ErrInvalidGroupPath
	}
	if filepath.Clean(p) != p {
		return ErrInvalidGroupPath
	}
	if strings.HasPrefix(p, "/sys/fs/cgroup") {
		return ErrInvalidGroupPath
	}
	return nil
}
