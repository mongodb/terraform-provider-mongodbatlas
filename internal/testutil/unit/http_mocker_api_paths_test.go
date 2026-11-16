package unit_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"sync/atomic"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiPathsParsing(t *testing.T) {
	specParts := unit.ReadAPISpecPaths()
	assert.GreaterOrEqual(t, len(specParts), 5)
	assert.Contains(t, specParts, "GET")
	processArgsPath := "/api/atlas/v2/groups/6746ceed6f62fc3c122a3e0e/clusters/test-acc-tf-c-7871793563057636102/processArgs"
	getPaths := specParts["GET"]
	found1, err := unit.FindNormalizedPath(processArgsPath, &getPaths)
	require.NoError(t, err)
	assert.Equal(t, "/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/processArgs", found1.Path)
	variables := found1.Variables(processArgsPath)
	assert.Equal(t, "6746ceed6f62fc3c122a3e0e", variables["groupId"])
	assert.Equal(t, "test-acc-tf-c-7871793563057636102", variables["clusterName"])
}

func TestDownloadOpenAPISpec(t *testing.T) {
	const attemptsBeforeSuccess = 3

	t.Run("downloads on first attempt", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			_, _ = w.Write([]byte("paths: {}"))
		}))
		defer server.Close()

		specPath := path.Join(t.TempDir(), "spec.yml")
		err := unit.DownloadOpenAPISpec(server.URL, specPath)
		require.NoError(t, err)
		assert.Equal(t, int32(1), attempts.Load())
		content, err := os.ReadFile(specPath)
		require.NoError(t, err)
		assert.Equal(t, "paths: {}", string(content))
	})

	t.Run("retries server errors", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if attempts.Add(1) < attemptsBeforeSuccess {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = w.Write([]byte("paths: {}"))
		}))
		defer server.Close()

		specPath := path.Join(t.TempDir(), "spec.yml")
		err := unit.DownloadOpenAPISpec(server.URL, specPath)
		require.NoError(t, err)
		assert.Equal(t, int32(attemptsBeforeSuccess), attempts.Load())
		content, err := os.ReadFile(specPath)
		require.NoError(t, err)
		assert.Equal(t, "paths: {}", string(content))
	})

	t.Run("fails after all attempts", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		specPath := path.Join(t.TempDir(), "spec.yml")
		err := unit.DownloadOpenAPISpec(server.URL, specPath)
		require.ErrorContains(t, err, "download failed after 3 attempts: unexpected status code 500")
		assert.Equal(t, int32(attemptsBeforeSuccess), attempts.Load())
	})
}
