package acc

import "fmt"

func TimeoutConfig(createTimeout, updateTimeout, deleteTimeout *string) string {
	createTimeoutConfig := ""
	updateTimeoutConfig := ""
	deleteTimeoutConfig := ""

	if createTimeout != nil {
		createTimeoutConfig = fmt.Sprintf(`
				create = %q
			`, *createTimeout)
	}
	if updateTimeout != nil {
		updateTimeoutConfig = fmt.Sprintf(`
			update = %q
		`, *updateTimeout)
	}
	if deleteTimeout != nil {
		deleteTimeoutConfig = fmt.Sprintf(`
			delete = %q
		`, *deleteTimeout)
	}
	timeoutConfig := "timeouts ="

	return fmt.Sprintf(`
		%[1]s {
			%[2]s
			%[3]s
			%[4]s
		}
	`, timeoutConfig, createTimeoutConfig, updateTimeoutConfig, deleteTimeoutConfig)
}

func ConfigRemove(resourceName string) string {
	return fmt.Sprintf(`
		removed {
			from = %s
			lifecycle {
				destroy = false
			}
		}
	`, resourceName)
}

func ConfigEmpty() string {
	return `
		# empty config to trigger delete
	`
}
