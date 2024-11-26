package acc

func ConvertAdvancedClusterToTPF(def string) string {
	if !IsTPFAdvancedCluster() {
		return def
	}
	return "invalid resource"
}
