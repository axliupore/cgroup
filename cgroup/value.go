package cgroup

import (
	"os"
	"path/filepath"
	"strconv"
)

type Value struct {
	filename string
	value    interface{}
}

func (v *Value) write(path string, perm os.FileMode) error {
	var data []byte
	switch t := v.value.(type) {
	case uint64:
		data = []byte(strconv.FormatUint(t, 10))
	case uint16:
		data = []byte(strconv.FormatUint(uint64(t), 10))
	case int64:
		data = []byte(strconv.FormatInt(t, 10))
	case []byte:
		data = t
	case string:
		data = []byte(t)
	case CPUMax:
		data = []byte(t)
	default:
		return ErrInvalidFormat
	}
	return os.WriteFile(
		filepath.Join(path, v.filename),
		data,
		perm,
	)
}

func writeValues(path string, values []Value) error {
	for _, v := range values {
		if err := v.write(path, filePerm); err != nil {
			return err
		}
	}
	return nil
}
