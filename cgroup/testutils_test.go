package cgroup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkFileContent(t *testing.T, path, filename, value string) {
	out, err := os.ReadFile(filepath.Join(path, filename))
	require.NoErrorf(t, err, "failed to read %s file", filename)
	assert.Equal(t, value, strings.TrimSpace(string(out)))
}
