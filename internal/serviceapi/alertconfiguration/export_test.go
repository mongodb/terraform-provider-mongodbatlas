package alertconfiguration

// UpgradeManualStateJSONForTest exposes the state upgrader's JSON rewrite for unit tests.
func UpgradeManualStateJSONForTest(raw []byte) ([]byte, error) {
	return upgradeManualStateJSON(raw)
}
