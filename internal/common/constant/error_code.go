package constant

// Atlas Admin API error codes shared across resources.
const (
	// ErrorCodeAPIKeyNotFound is the error code received when trying to delete an API key that does not exist.
	ErrorCodeAPIKeyNotFound = "API_KEY_NOT_FOUND" // #nosec G101 -- API error code, not a credential
)
