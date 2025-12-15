package constant

import "time"

const (
	DefaultTimeout = 3 * time.Hour
)

const (
	timeoutBaseDescription             = "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\". Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours)."
	timeoutDeleteAdditionalDescription = "Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs."
)

func TimeoutDescriptionCreateReadUpdate(defaultValue string) string {
	return timeoutBaseDescription + " Default: `" + defaultValue + "`."
}

func TimeoutDescriptionDelete(defaultValue string) string {
	return timeoutBaseDescription + " " + timeoutDeleteAdditionalDescription + " Default: `" + defaultValue + "`."
}
