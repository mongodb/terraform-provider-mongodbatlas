package replay

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	hoverfly "github.com/SpectoLabs/hoverfly/core"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

const simulationDir = "../../../simulations/"

func IsInCaptureMode() bool {
	return os.Getenv("REPLAY_MODE") == "capture"
}

func IsInSimulateMode() bool {
	return os.Getenv("REPLAY_MODE") == "simulate"
}

func SetupReplayProxy(t *testing.T) (proxyPort *int, teardown func(t *testing.T)) {
	t.Helper()
	if IsInCaptureMode() {
		return setupCaptureMode()
	}
	if IsInSimulateMode() {
		return setupSimulateMode(t)
	}

	log.Printf("No replay mode was configured")
	return nil, func(t *testing.T) {
		t.Helper()
	}
}

func setupSimulateMode(t *testing.T) (proxyPort *int, teardown func(t *testing.T)) {
	t.Helper()
	port, hv := newHoverflyInstance()

	fileName := fmt.Sprintf("%s%s.json", simulationDir, t.Name())
	if err := hv.ImportFromDisk(fileName); err != nil {
		log.Fatalf("Failed to import simulation for test: %v", err)
	}

	if err := hv.SetMode("simulate"); err != nil {
		log.Fatalf("Failed to set Hoverfly mode: %v", err)
	}

	return &port, teardownSimulate(hv)
}

func setupCaptureMode() (proxyPort *int, teardown func(t *testing.T)) {
	port, hv := newHoverflyInstance()

	if err := hv.SetModeWithArguments(v2.ModeView{Mode: "capture", Arguments: v2.ModeArgumentsView{Stateful: true}}); err != nil {
		log.Fatalf("Failed to set Hoverfly mode: %v", err)
	}
	return &port, teardownCapture(hv)
}

func newHoverflyInstance() (int, *hoverfly.Hoverfly) {
	var min int64 = 1024
	var max int64 = 65536
	diff := max - min
	nBig, err := rand.Int(rand.Reader, big.NewInt(diff+1))
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	proxyPort := int(nBig.Int64() + min)

	settings := hoverfly.InitSettings()
	settings.ProxyPort = fmt.Sprintf("%d", proxyPort)
	hv := hoverfly.NewHoverflyWithConfiguration(settings)

	if err := hv.StartProxy(); err != nil {
		log.Fatalf("Failed to start Hoverfly: %v", err)
	}
	return proxyPort, hv
}

func teardownSimulate(hv *hoverfly.Hoverfly) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		hv.StopProxy()
	}
}

func teardownCapture(hv *hoverfly.Hoverfly) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		data, err := hv.GetSimulation()
		if err != nil {
			log.Fatalf("Failed to obtain simulation result: %v", err)
		}
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Fatalf("Error serializing to JSON: %v", err)
		}

		if err := createFileInSimulationDir(jsonData, fmt.Sprintf("%s.json", t.Name())); err != nil {
			log.Fatalf("Error storing file: %v", err)
		}
		hv.StopProxy()
	}
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
