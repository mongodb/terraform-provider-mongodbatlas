package replay

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	hoverfly "github.com/SpectoLabs/hoverfly/core"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const resultFilePrefix = "../../../simulations/"

func IsInCaptureMode() bool {
	return os.Getenv("REPLAY_MODE") == "capture"
}

func IsInSimulateMode() bool {
	return os.Getenv("REPLAY_MODE") == "simulate"
}

func SetupReplayProxy(t *testing.T) (*int, func(t *testing.T)) {
	if IsInCaptureMode() {
		proxyPort, hv := newHoverflyInstance()

		if err := hv.SetModeWithArguments(v2.ModeView{Mode: "capture", Arguments: v2.ModeArgumentsView{Stateful: true}}); err != nil {
			log.Fatalf("Failed to set Hoverfly mode: %v", err)
		}
		return &proxyPort, teardownCapture(hv)
	}
	if IsInSimulateMode() {
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

	log.Printf("No replay mode was configured")
	return nil, func(t *testing.T) {}
}

func newHoverflyInstance() (int, *hoverfly.Hoverfly) {
	proxyPort := acctest.RandIntRange(1024, 65536)
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
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Fatalf("Error serializing to JSON: %v", err)
		}

		if err := createFileInSimulationDir(jsonData, t.Name()); err != nil {
			log.Fatalf("Error storing file: %v", err)
		}
		hv.StopProxy()
	}
}

func createFileInSimulationDir(jsonData []byte, fileName string) error {
	filePath := fmt.Sprintf("%s%s.json", resultFilePrefix, fileName)
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}
