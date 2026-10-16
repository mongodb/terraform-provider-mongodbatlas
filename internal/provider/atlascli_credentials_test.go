package provider_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

func TestGetAtlasCLICredentials(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on a POSIX shell script named `atlas`")
	}

	writeFakeAtlas := func(t *testing.T, body string) {
		t.Helper()
		dir := t.TempDir()
		path := filepath.Join(dir, "atlas")
		require.NoError(t, os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0o700)) //nolint:gosec // fake test executable must be executable
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	}

	t.Run("returns access token when logged in", func(t *testing.T) {
		writeFakeAtlas(t, `echo "atlas-cli-token"`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context())
		require.NoError(t, err)
		require.NotNil(t, creds)
		assert.Equal(t, "atlas-cli-token", creds.AccessToken)
	})

	t.Run("returns nil when not logged in", func(t *testing.T) {
		writeFakeAtlas(t, `echo "no access token found" 1>&2; exit 1`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context())
		require.NoError(t, err)
		assert.Nil(t, creds)
	})

	t.Run("returns nil when output is empty", func(t *testing.T) {
		writeFakeAtlas(t, `exit 0`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context())
		require.NoError(t, err)
		assert.Nil(t, creds)
	})

	t.Run("returns nil when atlas is not installed", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context())
		require.NoError(t, err)
		assert.Nil(t, creds)
	})
}
