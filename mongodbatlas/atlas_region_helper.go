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
	// Europe
	EuWest1    AtlasRegion = "EU_WEST_1"
	EuCentral1 AtlasRegion = "EU_CENTRAL_1"
	EuNorth1   AtlasRegion = "EU_NORTH_1"
	EuWest2    AtlasRegion = "EU_WEST_2"
	EuWest3    AtlasRegion = "EU_WEST_3"
	EuSouth1   AtlasRegion = "EU_SOUTH_1"
	EuCentral2 AtlasRegion = "EU_CENTRAL_2"
	EuSouth2   AtlasRegion = "EU_SOUTH_2"
	// Asia Pacific
	ApSoutheast1 AtlasRegion = "AP_SOUTHEAST_1"
	ApSoutheast2 AtlasRegion = "AP_SOUTHEAST_2"
	ApSoutheast3 AtlasRegion = "AP_SOUTHEAST_3"
	ApSouth1     AtlasRegion = "AP_SOUTH_1"
	ApEast1      AtlasRegion = "AP_EAST_1"
	ApNortheast1 AtlasRegion = "AP_NORTHEAST_1"
	ApNortheast2 AtlasRegion = "AP_NORTHEAST_2"
	ApNortheast3 AtlasRegion = "AP_NORTHEAST_3"
	ApSouth2     AtlasRegion = "AP_SOUTH_2"
	ApSoutheast4 AtlasRegion = "AP_SOUTHEAST_4"
	// Middle East and Africa
	MeSouth1   AtlasRegion = "ME_SOUTH_1"
	AfSouth1   AtlasRegion = "AF_SOUTH_1"
	MeCentral1 AtlasRegion = "ME_CENTRAL_1"

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
	// Asia Pacific
	AsiaEast           AtlasRegion = "ASIA_EAST"
	AsiaSouthEast      AtlasRegion = "ASIA_SOUTH_EAST"
	AustraliaCentral   AtlasRegion = "AUSTRALIA_CENTRAL"
	AustraliaCentral2  AtlasRegion = "AUSTRALIA_CENTRAL_2"
	AustraliaEast      AtlasRegion = "AUSTRALIA_EAST"
	AustraliaSouthEast AtlasRegion = "AUSTRALIA_SOUTH_EAST"
	IndiaCentral       AtlasRegion = "INDIA_CENTRAL"
	IndiaSouth         AtlasRegion = "INDIA_SOUTH"
	IndiaWest          AtlasRegion = "INDIA_WEST"
	JapanEast          AtlasRegion = "JAPAN_EAST"
	JapanWest          AtlasRegion = "JAPAN_WEST"
	KoreaCentral       AtlasRegion = "KOREA_CENTRAL"
	KoreaSouth         AtlasRegion = "KOREA_SOUTH"
	// Europe
	EuropeNorth        AtlasRegion = "EUROPE_NORTH"
	EuropeWest         AtlasRegion = "EUROPE_WEST"
	UkSouth            AtlasRegion = "UK_SOUTH"
	UkWest             AtlasRegion = "UK_WEST"
	FranceCentral      AtlasRegion = "FRANCE_CENTRAL"
	FranceSouth        AtlasRegion = "FRANCE_SOUTH"
	GermanyWestCentral AtlasRegion = "GERMANY_WEST_CENTRAL"
	GermanyNorth       AtlasRegion = "GERMANY_NORTH"
	SwitzerlandNorth   AtlasRegion = "SWITZERLAND_NORTH"
	SwitzerlandWest    AtlasRegion = "SWITZERLAND_WEST"
	NorwayEast         AtlasRegion = "NORWAY_EAST"
	NorwayWest         AtlasRegion = "NORWAY_WEST"
	SwedenCentral      AtlasRegion = "SWEDEN_CENTRAL"
	SwedenSouth        AtlasRegion = "SWEDEN_SOUTH"
	// Middle East and Africa
	SouthAfricaNorth AtlasRegion = "SOUTH_AFRICA_NORTH"
	SouthAfricaWest  AtlasRegion = "SOUTH_AFRICA_WEST"
	UaeNorth         AtlasRegion = "UAE_NORTH"
	UaeCentral       AtlasRegion = "UAE_CENTRAL"
	QatarCentral     AtlasRegion = "QATAR_CENTRAL"

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
	// Asia Pacific
	EasternAsiaPacific      AtlasRegion = "EASTERN_ASIA_PACIFIC"
	AsiaEast2               AtlasRegion = "ASIA_EAST_2"
	NortheasternAsiaPacific AtlasRegion = "NORTHEASTERN_ASIA_PACIFIC"
	AsiaNortheast2          AtlasRegion = "ASIA_NORTHEAST_2"
	AsiaNortheast3          AtlasRegion = "ASIA_NORTHEAST_3"
	SoutheasternAsiaPacific AtlasRegion = "SOUTHEASTERN_ASIA_PACIFIC"
	AsiaSouth1              AtlasRegion = "ASIA_SOUTH_1"
	AsiaSouth2              AtlasRegion = "ASIA_SOUTH_2"
	AustraliaSoutheast1     AtlasRegion = "AUSTRALIA_SOUTHEAST_1"
	AustraliaSoutheast2     AtlasRegion = "AUSTRALIA_SOUTHEAST_2"
	AsiaSoutheast2          AtlasRegion = "ASIA_SOUTHEAST_2"
	// Europe
	WesternEurope    AtlasRegion = "WESTERN_EUROPE"
	EuropeNorth1     AtlasRegion = "EUROPE_NORTH_1"
	EuropeWest2      AtlasRegion = "EUROPE_WEST_2"
	EuropeWest3      AtlasRegion = "EUROPE_WEST_3"
	EuropeWest4      AtlasRegion = "EUROPE_WEST_4"
	EuropeWest6      AtlasRegion = "EUROPE_WEST_6"
	EuropeCentral2   AtlasRegion = "EUROPE_CENTRAL_2"
	EuropeWest8      AtlasRegion = "EUROPE_WEST_8"
	EuropeWest9      AtlasRegion = "EUROPE_WEST_9"
	EuropeWest12     AtlasRegion = "EUROPE_WEST_12"
	EuropeSouthwest1 AtlasRegion = "EUROPE_SOUTHWEST_1"
	// Middle East and Africa
	MiddleEastWest1    AtlasRegion = "MIDDLE_EAST_WEST_1"
	MiddleEastCentral1 AtlasRegion = "MIDDLE_EAST_CENTRAL_1"
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
		// Asia Pacific
		"ap-southeast-1": ApSoutheast1,
		"ap-southeast-2": ApSoutheast2,
		"ap-southeast-3": ApSoutheast3,
		"ap-south-1":     ApSouth1,
		"ap-east-1":      ApEast1,
		"ap-northeast-1": ApNortheast1,
		"ap-northeast-2": ApNortheast2,
		"ap-northeast-3": ApNortheast3,
		"ap-south-2":     ApSouth2,
		"ap-southeast-4": ApSoutheast4,
		// Europe
		"eu-west-1":    EuWest1,
		"eu-central-1": EuCentral1,
		"eu-north-1":   EuNorth1,
		"eu-west-2":    EuWest2,
		"eu-west-3":    EuWest3,
		"eu-south-1":   EuSouth1,
		"eu-central-2": EuCentral2,
		"eu-south-2":   EuSouth2,
		// Middle East and Africa
		"me-south-1":   MeSouth1,
		"af-south-1":   AfSouth1,
		"me-central-1": MeCentral1,
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
		// Europe
		"northeurope":        EuropeNorth,
		"westeurope":         EuropeWest,
		"uksouth":            UkSouth,
		"ukwest":             UkWest,
		"francecentral":      FranceCentral,
		"francesouth":        FranceSouth,
		"germanywestcentral": GermanyWestCentral,
		"germanynorth":       GermanyNorth,
		"switzerlandnorth":   SwitzerlandNorth,
		"switzerlandwest":    SwitzerlandWest,
		"norwayeast":         NorwayEast,
		"norwaywest":         NorwayWest,
		"swedencentral":      SwedenCentral,
		"swedensouth":        SwedenSouth,
		// Asia Pacific
		"eastasia":           AsiaEast,
		"southeastasia":      AsiaSouthEast,
		"australiacentral":   AustraliaCentral,
		"australiacentral2":  AustraliaCentral2,
		"australiaeast":      AustraliaEast,
		"australiasoutheast": AustraliaSouthEast,
		"centralindia":       IndiaCentral,
		"southindia":         IndiaSouth,
		"westindia":          IndiaWest,
		"japaneast":          JapanEast,
		"japanwest":          JapanWest,
		"koreacentral":       KoreaCentral,
		"koreasouth":         KoreaSouth,
		// Middle East and Africa
		"southafricanorth": SouthAfricaNorth,
		"southafricawest":  SouthAfricaWest,
		"uaenorth":         UaeNorth,
		"uaecentral":       UaeCentral,
		"qatarcentral":     QatarCentral,
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
		// Asia Pacific
		"asia-east1":           EasternAsiaPacific,
		"asia-east2":           AsiaEast2,
		"asia-northeast1":      NortheasternAsiaPacific,
		"asia-northeast2":      AsiaNortheast2,
		"asia-northeast3":      AsiaNortheast3,
		"asia-southeast1":      SoutheasternAsiaPacific,
		"asia-south1":          AsiaSouth1,
		"asia-south2":          AsiaSouth2,
		"australia-southeast1": AustraliaSoutheast1,
		"australia-southeast2": AustraliaSoutheast2,
		"asia-southeast2":      AsiaSoutheast2,
		// Europe
		"europe-west1":      WesternEurope,
		"europe-north1":     EuropeNorth1,
		"europe-west2":      EuropeWest2,
		"europe-west3":      EuropeWest3,
		"europe-west4":      EuropeWest4,
		"europe-west6":      EuropeWest6,
		"europe-central2":   EuropeCentral2,
		"europe-west8":      EuropeWest8,
		"europe-west9":      EuropeWest9,
		"europe-west12":     EuropeWest12,
		"europe-southwest1": EuropeSouthwest1,
		// Middle East
		"me-west1":    MiddleEastWest1,
		"me-central1": MiddleEastCentral1,
	}
}
