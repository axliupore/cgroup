package cgroup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyGroupPath(t *testing.T) {
	valids := map[string]bool{
		"/":                                true,
		"":                                 false,
		"/axliupore":                       true,
		"/axliupore/xiaoyu":                true,
		"/sys/fs/cgroup/axliupore":         false,
		"/sys/fs/cgroup/unified/axliupore": false,
		"axliupore":                        false,
		"/axliupore/../xiaoyu":             false,
	}
	for s, valid := range valids {
		err := VerifyGroupPath(s)
		if valid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
