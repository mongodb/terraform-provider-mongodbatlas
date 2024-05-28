package replay

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const simulationDir = "../../../simulations/"

func IsInCaptureMode() bool {
	return os.Getenv("REPLAY_MODE") == "capture"
}

func IsInSimulateMode() bool {
	return os.Getenv("REPLAY_MODE") == "simulate"
}

func SetupReplayProxy(t *testing.T) *int {
	t.Helper()
	if IsInCaptureMode() {
		return setupCaptureMode(t)
	}
	if IsInSimulateMode() {
		return setupSimulateMode(t)
	}

	log.Printf("No replay mode was configured")
	return nil
}

func setupSimulateMode(t *testing.T) *int {
	t.Helper()
	port := randomPortNumber()
	adminPort := port + 1
	simulationFilePath := fmt.Sprintf("%s%s.json", simulationDir, t.Name())
	cmd := exec.Command("../../../scripts/hoverfly-import-and-simulate.sh", fmt.Sprintf("%d", port), fmt.Sprintf("%d", adminPort), simulationFilePath) //nolint:gosec // inputs are fully controlled within function
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start hoverfly in simulate mode: %s", err)
	}
	t.Cleanup(func() { cleanupSimulateMode(t, port) })
	return &port
}

func setupCaptureMode(t *testing.T) (proxyPort *int) {
	t.Helper()
	port := randomPortNumber()
	adminPort := port + 1
	cmd := exec.Command("../../../scripts/hoverfly-capture-mode.sh", fmt.Sprintf("%d", port), fmt.Sprintf("%d", adminPort)) //nolint:gosec // inputs are fully controlled within function
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start hoverfly in capture mode: %s", err)
	}
	t.Cleanup(func() { cleanupCaptureMode(t, port) })
	return &port
}

func cleanupSimulateMode(t *testing.T, port int) {
	t.Helper()
	cmd := exec.Command("../../../scripts/hoverfly-end-simulation.sh", fmt.Sprintf("%d", port)) //nolint:gosec // inputs are fully controlled within function
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to end hoverfly simulation: %s", err)
	}
}

func cleanupCaptureMode(t *testing.T, port int) {
	t.Helper()
	simulationFilePath := fmt.Sprintf("%s%s.json", simulationDir, t.Name())
	cmd := exec.Command("../../../scripts/hoverfly-export-simulation.sh", fmt.Sprintf("%d", port), simulationFilePath) //nolint:gosec // inputs are fully controlled within function
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stop and export: %s", err)
	}
}

func randomPortNumber() int {
	var minPort int64 = 1024
	var maxPort int64 = 65536
	diff := maxPort - minPort
	nBig, err := rand.Int(rand.Reader, big.NewInt(diff+1))
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	proxyPort := int(nBig.Int64() + minPort)
	return proxyPort
}

func createFileInSimulationDir(jsonData []byte, fileName string) error {
	filePath := filePathInSimulationDir(fileName)
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0o600)
}

func filePathInSimulationDir(fileName string) string {
	return fmt.Sprintf("%s%s", simulationDir, fileName)
}
