package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRunMain is a quick test that runs main() function
func TestRunMain(t *testing.T) {

	// This is a simple integration test that runs main
	// We're just checking it runs without panicking
	assert.NotPanics(t, func() {
		main()
	}, "main() should execute without panicking")

	// Verify spec file exists
	_, err := os.Stat(specFilePath)
	assert.NoError(t, err, "Expected spec file to exist after running main")
}
