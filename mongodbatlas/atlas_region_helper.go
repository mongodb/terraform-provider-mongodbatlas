package mongodbatlas

import (
	"fmt"
	"regexp"
	"strings"
)

type AtlasRegion string

const (
	UsEast1    AtlasRegion = "US_EAST_1"
	UsWest2    AtlasRegion = "US_WEST_2"
	CaCentral1 AtlasRegion = "CA_CENTRAL_1"
	UsEast2    AtlasRegion = "US_EAST_2"
	UsWest1    AtlasRegion = "US_WEST_1"
	SaEast1    AtlasRegion = "SA_EAST_1"

	// Azure
	UsCentral
	UsEast
	UsNorthCentral
	UsWest

	// GCP
	CentralUs
	EasternUs
	UsEast4
	USEast5
	NorthAmericaNortheast1
	NorthAmericaNortheast2
)

func IsAtlasRegion(regionName string) bool {
	var regex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return regex.MatchString(regionName)
}

func GetAtlasRegion(cloudProviderName, backingProviderName, regionName string) (string, error) {
	var atlasRegionName AtlasRegion
	var exists bool

	if IsAtlasRegion(regionName) {
		return regionName, nil
	}

	regionsMap := cloudProviderToAtlasRegionsMap()

	if cloudProviderName == "TENANT" {
		atlasRegionName, exists = regionsMap[backingProviderName][strings.ToLower(regionName)]
	} else {
		atlasRegionName, exists = regionsMap[cloudProviderName][strings.ToLower(regionName)]
	}

	if !exists {
		return "", fmt.Errorf("No Atlas region exists for cloud provider: %s, region name: %s", cloudProviderName, regionName)
	}

	return string(atlasRegionName), nil
}

func cloudProviderToAtlasRegionsMap() map[string]map[string]AtlasRegion {
	return map[string]map[string]AtlasRegion{
		"AWS":   awsToAtlasRegionsMap(),
		"AZURE": azureToAtlasRegionsMap(),
		"GCP":   gcpToAtlasRegionsMap(),
	}
}

func awsToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"us-east-1":    UsEast1,
		"us-west-2":    UsWest2,
		"ca-central-1": CaCentral1,
		"us-east-2":    UsEast2,
		"us-west-1":    UsWest1,
		"sa-east-1":    SaEast1,
	}
}

func azureToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"centralus":      UsCentral,
		"eastus":         UsEast,
		"eastus2":        UsEast2,
		"northcentralus": UsNorthCentral,
		"westus":         UsWest,
		"westus2":        UsWest2,
	}
}

func gcpToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"us-central1":             CentralUs,
		"us-east1":                EasternUs,
		"us-east4":                UsEast4,
		"us-east5":                USEast5,
		"northamerica-northeast1": NorthAmericaNortheast1,
		"northamerica-northeast2": NorthAmericaNortheast2,
	}
}
