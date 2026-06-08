package constant

import "time"

const (
	DefaultTimeout              = 3 * time.Hour
	DefaultTimeoutDocumentation = "3h"
)

const (
	timeoutCommonDescription           = "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\". Valid time units are \"s\" (seconds), \"m\" (minutes), and \"h\" (hours)."
	timeoutDeleteAdditionalDescription = "Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs."
	timeoutReadAdditionalDescription   = "Read operations occur during any refresh or planning operation when refresh is enabled."
)

func TimeoutDescriptionCreateUpdate(defaultValue string) string {
	return timeoutCommonDescription + " Default: `" + defaultValue + "`."
}

func TimeoutDescriptionRead(defaultValue string) string {
	return timeoutCommonDescription + " " + timeoutReadAdditionalDescription + " Default: `" + defaultValue + "`."
}

func TimeoutDescriptionDelete(defaultValue string) string {
	return timeoutCommonDescription + " " + timeoutDeleteAdditionalDescription + " Default: `" + defaultValue + "`."
}
