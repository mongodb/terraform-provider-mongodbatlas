package mongodbatlas

import (
	"fmt"
	"regexp"
)

type AtlasRegion string

const (
	US_EAST_1    AtlasRegion = "US_EAST_1"
	US_WEST_2    AtlasRegion = "US_WEST_2"
	CA_CENTRAL_1 AtlasRegion = "CA_CENTRAL_1"
	US_EAST_2    AtlasRegion = "US_EAST_2"
	US_WEST_1    AtlasRegion = "US_WEST_1"
	SA_EAST_1    AtlasRegion = "SA_EAST_1"

	// Azure
	US_CENTRAL
	US_EAST
	US_NORTH_CENTRAL
	US_WEST

	// GCP
	CENTRAL_US
	EASTERN_US
	US_EAST_4
	US_EAST_5
	NORTH_AMERICA_NORTHEAST_1
	NORTH_AMERICA_NORTHEAST_2
)

func IsAtlasRegion(regionName string) bool {
	var regex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return regex.MatchString(regionName)
}

func GetAtlasRegion(cloudProvider, regionName string) (string, error) {
	if IsAtlasRegion(regionName) {
		return regionName, nil
	}

	regionsMap := cloudProviderToAtlasRegionsMap()
	atlasRegionName, exists := regionsMap[cloudProvider][regionName]
	if !exists {
		return "", fmt.Errorf("No Atlas region exists for cloud provider: %s, region name: %s", cloudProvider, regionName)
	}

	return string(atlasRegionName), nil
}

// func AtlasRegionForCloudProviderRegion(cloudProvider, regionName string) (AtlasRegion, error) {
// 	regionsMap := cloudProviderToAtlasRegionsMap()
// 	atlasRegionName, exists := regionsMap[cloudProvider][regionName]
// 	if !exists {

// 	}

// 	return atlasRegionName, nil
// }

/*
AWS = us-east-1 | US_EAST_1

	us-west-2 | US_WEST_2
	.

AZURE = centralus | US_CENTRAL

	    eastus | US_EAST
		.

GCP = ...
*/
func cloudProviderToAtlasRegionsMap() map[string]map[string]AtlasRegion {
	return map[string]map[string]AtlasRegion{
		"AWS":   awsToAtlasRegionsMap(),
		"AZURE": azureToAtlasRegionsMap(),
		"GCP":   gcpToAtlasRegionsMap(),
	}
}

func awsToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"us-east-1":    US_EAST_1,
		"us-west-2":    US_WEST_2,
		"ca-central-1": CA_CENTRAL_1,
		"us-east-2":    US_EAST_2,
		"us-west-1":    US_WEST_1,
		"sa-east-1":    SA_EAST_1,
	}
}

func azureToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"centralus":      US_CENTRAL,
		"eastus":         US_EAST,
		"eastus2":        US_EAST_2,
		"northcentralus": US_NORTH_CENTRAL,
		"westus":         US_WEST,
		"westus2":        US_WEST_2,
	}
}

func gcpToAtlasRegionsMap() map[string]AtlasRegion {
	return map[string]AtlasRegion{
		"us-central1":             CENTRAL_US,
		"us-east1":                EASTERN_US,
		"us-east4":                US_EAST_4,
		"us-east5":                US_EAST_5,
		"northamerica-northeast1": NORTH_AMERICA_NORTHEAST_1,
		"northamerica-northeast2": NORTH_AMERICA_NORTHEAST_2,
	}
}
