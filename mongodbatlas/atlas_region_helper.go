package mongodbatlas

import (
	"fmt"
	"regexp"
	"strings"
)

type AtlasRegion string

const (
	// AWS:
	// Americas
	UsEast1    AtlasRegion = "US_EAST_1"
	UsWest2    AtlasRegion = "US_WEST_2"
	CaCentral1 AtlasRegion = "CA_CENTRAL_1"
	UsEast2    AtlasRegion = "US_EAST_2"
	UsWest1    AtlasRegion = "US_WEST_1"
	SaEast1    AtlasRegion = "SA_EAST_1"

	// Azure:
	// Americas
	UsCentral       AtlasRegion = "US_CENTRAL"
	UsEast          AtlasRegion = "US_EAST"
	UsNorthCentral  AtlasRegion = "US_NORTH_CENTRAL"
	UsWest          AtlasRegion = "US_WEST"
	UsWest3         AtlasRegion = "US_WEST_3"
	UsWestCentral   AtlasRegion = "US_WEST_CENTRAL"
	UsSouthCentral  AtlasRegion = "US_SOUTH_CENTRAL"
	BrazilSouth     AtlasRegion = "BRAZIL_SOUTH"
	BrazilSoutheast AtlasRegion = "BRAZIL_SOUTHEAST"
	CanadaEast      AtlasRegion = "CANADA_EAST"
	CanadaCentral   AtlasRegion = "CANADA_CENTRAL"

	// GCP:
	// Americas
	CentralUs              AtlasRegion = "CENTRAL_US"
	EasternUs              AtlasRegion = "EASTERN_US"
	UsEast4                AtlasRegion = "US_EAST_4"
	USEast5                AtlasRegion = "US_EAST_5"
	NorthAmericaNortheast1 AtlasRegion = "NORTH_AMERICA_NORTHEAST_1"
	NorthAmericaNortheast2 AtlasRegion = "NORTH_AMERICA_NORTHEAST_2"
	SouthAmericaEast1      AtlasRegion = "SOUTH_AMERICA_EAST_1"
	SouthAmericaWest1      AtlasRegion = "SOUTH_AMERICA_WEST_1"
	WesternUs              AtlasRegion = "WESTERN_US"
	UsWest4                AtlasRegion = "US_WEST_4"
	UsSouth1               AtlasRegion = "US_SOUTH_1"
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
		return "", fmt.Errorf("no Atlas region exists for cloud provider: %s, region name: %s", cloudProviderName, regionName)
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
		"centralus":       UsCentral,
		"eastus":          UsEast,
		"eastus2":         UsEast2,
		"northcentralus":  UsNorthCentral,
		"westus":          UsWest,
		"westus2":         UsWest2,
		"westus3":         UsCentral,
		"westcentralus":   UsEast,
		"southcentralus":  UsEast2,
		"brazilsouth":     UsNorthCentral,
		"brazilsoutheast": UsWest,
		"canadaeast":      UsWest2,
		"canadacentral":   UsWest2,
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
		"southamerica-east1":      SouthAmericaEast1,
		"southamerica-west1":      SouthAmericaWest1,
		"us-west1":                UsWest1,
		"us-west2":                UsWest2,
		"us-west3":                UsWest3,
		"us-west4":                UsWest4,
		"us-south1":               UsSouth1,
	}
}
