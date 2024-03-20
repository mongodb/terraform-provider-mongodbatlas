package replay

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	hoverfly "github.com/SpectoLabs/hoverfly/core"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

const resultFilePrefix = "../../../simulations/"

func SetupReplayProxy(t *testing.T) (*int, func(t *testing.T)) {

	configuredMode := os.Getenv("PROXY_MODE")

	if configuredMode == "capture" {
		proxyPort, hv := newHoverflyInstance()

		if err := hv.SetModeWithArguments(v2.ModeView{Mode: "capture", Arguments: v2.ModeArgumentsView{Stateful: true}}); err != nil {
			log.Fatalf("Failed to set Hoverfly mode: %v", err)
		}
		return &proxyPort, teardownCapture(hv)
	}
	if configuredMode == "simulate" {
		proxyPort, hv := newHoverflyInstance()

		fileName := fmt.Sprintf("%s%s.json", resultFilePrefix, t.Name())
		if err := hv.ImportFromDisk(fileName); err != nil {
			log.Fatalf("Failed to import simulation for test: %v", err)
		}

		if err := hv.SetMode("simulate"); err != nil {
			log.Fatalf("Failed to set Hoverfly mode: %v", err)
		}

		return &proxyPort, teardownSimulate(hv)
	}

	log.Printf("No proxy valid proxy mode was configured: %s", configuredMode)
	return nil, func(t *testing.T) {}
}

func newHoverflyInstance() (int, *hoverfly.Hoverfly) {
	proxyPort := rand.Intn(65536)
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
		hv.StopProxy()
	}
}

func teardownCapture(hv *hoverfly.Hoverfly) func(t *testing.T) {
	return func(t *testing.T) {
		data, err := hv.GetSimulation()
		if err != nil {
			log.Fatalf("Failed to obtain simulation result: %v", err)
		}
		fileName := fmt.Sprintf("%s%s.json", resultFilePrefix, t.Name())
		serializeSimulationToFile(data, fileName)
		hv.StopProxy()
	}
}

func serializeSimulationToFile(simulation v2.SimulationViewV5, filePath string) {
	jsonData, err := json.MarshalIndent(simulation, "", "    ")
	if err != nil {
		log.Fatalf("Error serializing to JSON: %v", err)
	}

	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatalf("Error creating directories: %v", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}
}
