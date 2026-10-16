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

// writeFakeAtlas installs a fake `atlas` executable on PATH. Every invocation is
// appended to the returned log file so tests can assert the arguments used.
func writeFakeAtlas(t *testing.T, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("test relies on a POSIX shell script named `atlas`")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "atlas")
	logFile := filepath.Join(dir, "calls.log")
	script := "#!/bin/sh\necho \"$@\" >> \"$ATLAS_LOG\"\n" + body + "\n"
	require.NoError(t, os.WriteFile(path, []byte(script), 0o700)) //nolint:gosec // fake test executable must be executable
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("ATLAS_LOG", logFile)
	return logFile
}

func TestGetAtlasCLICredentials(t *testing.T) {
	t.Run("returns token with base_url and gov from the CLI profile", func(t *testing.T) {
		writeFakeAtlas(t, `case "$1 $2" in
  "auth token") echo "atlas-cli-token" ;;
  "config get") echo '{"ops_manager_url":"https://cloud-dev.mongodbgov.com","service":"cloudgov"}' ;;
  *) exit 1 ;;
esac`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "default")
		require.NoError(t, err)
		require.NotNil(t, creds)
		assert.Equal(t, "atlas-cli-token", creds.AccessToken)
		assert.Equal(t, "https://cloud-dev.mongodbgov.com", creds.BaseURL)
		assert.True(t, creds.IsMongodbGovCloud)
	})

	t.Run("returns token without base_url when config get fails", func(t *testing.T) {
		writeFakeAtlas(t, `case "$1 $2" in
  "auth token") echo "atlas-cli-token" ;;
  *) exit 1 ;;
esac`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "default")
		require.NoError(t, err)
		require.NotNil(t, creds)
		assert.Equal(t, "atlas-cli-token", creds.AccessToken)
		assert.Empty(t, creds.BaseURL)
		assert.False(t, creds.IsMongodbGovCloud)
	})

	t.Run("returns nil when not logged in", func(t *testing.T) {
		writeFakeAtlas(t, `exit 1`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "default")
		require.NoError(t, err)
		assert.Nil(t, creds)
	})

	t.Run("returns nil when output is empty", func(t *testing.T) {
		writeFakeAtlas(t, `exit 0`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "default")
		require.NoError(t, err)
		assert.Nil(t, creds)
	})

	t.Run("returns nil when atlas is not installed", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "default")
		require.NoError(t, err)
		assert.Nil(t, creds)
	})

	t.Run("passes the given profile to both commands", func(t *testing.T) {
		logFile := writeFakeAtlas(t, `case "$1 $2" in
  "auth token") echo "atlas-cli-token" ;;
  "config get") echo '{"ops_manager_url":"https://cloud-dev.mongodb.com/","service":"cloud"}' ;;
  *) exit 1 ;;
esac`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "myprofile")
		require.NoError(t, err)
		require.NotNil(t, creds)
		calls, err := os.ReadFile(logFile)
		require.NoError(t, err)
		assert.Contains(t, string(calls), "auth token --profile myprofile")
		assert.Contains(t, string(calls), "config get myprofile -o json")
	})

	t.Run("lowercases a mixed-case profile for both commands", func(t *testing.T) {
		logFile := writeFakeAtlas(t, `case "$1 $2" in
  "auth token") echo "atlas-cli-token" ;;
  "config get") echo '{"ops_manager_url":"https://cloud-dev.mongodb.com/","service":"cloud"}' ;;
  *) exit 1 ;;
esac`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "Dev")
		require.NoError(t, err)
		require.NotNil(t, creds)
		assert.Equal(t, "https://cloud-dev.mongodb.com/", creds.BaseURL)
		calls, err := os.ReadFile(logFile)
		require.NoError(t, err)
		assert.Contains(t, string(calls), "auth token --profile dev")
		assert.Contains(t, string(calls), "config get dev -o json")
	})

	t.Run("resolves the default profile when none is given", func(t *testing.T) {
		logFile := writeFakeAtlas(t, `case "$1 $2" in
  "config list") echo '["only"]' ;;
  "auth token") echo "atlas-cli-token" ;;
  "config get") echo '{}' ;;
  *) exit 1 ;;
esac`)
		creds, err := provider.GetAtlasCLICredentialsForTest(t.Context(), "")
		require.NoError(t, err)
		require.NotNil(t, creds)
		calls, err := os.ReadFile(logFile)
		require.NoError(t, err)
		assert.Contains(t, string(calls), "auth token --profile only")
		assert.Contains(t, string(calls), "config get only -o json")
	})
}

func TestDefaultAtlasCLIProfile(t *testing.T) {
	t.Run("uses the only profile when a single one exists", func(t *testing.T) {
		writeFakeAtlas(t, `case "$1 $2" in "config list") echo '["only"]' ;; *) exit 1 ;; esac`)
		assert.Equal(t, "only", provider.DefaultAtlasCLIProfileForTest(t.Context()))
	})

	t.Run("falls back to default with multiple profiles", func(t *testing.T) {
		writeFakeAtlas(t, `case "$1 $2" in "config list") echo '["a","b"]' ;; *) exit 1 ;; esac`)
		assert.Equal(t, "default", provider.DefaultAtlasCLIProfileForTest(t.Context()))
	})

	t.Run("falls back to default when config list fails", func(t *testing.T) {
		writeFakeAtlas(t, `exit 1`)
		assert.Equal(t, "default", provider.DefaultAtlasCLIProfileForTest(t.Context()))
	})

	t.Run("falls back to default on malformed config list output", func(t *testing.T) {
		writeFakeAtlas(t, `case "$1 $2" in "config list") echo 'not-json' ;; *) exit 1 ;; esac`)
		assert.Equal(t, "default", provider.DefaultAtlasCLIProfileForTest(t.Context()))
	})
}
