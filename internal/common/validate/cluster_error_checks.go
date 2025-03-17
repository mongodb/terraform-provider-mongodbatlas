package validate

import admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"

func ErrorClusterIsAsymmetrics(err error) bool {
	return err != nil && admin20240530.IsErrorCode(err, "ASYMMETRIC_SHARD_UNSUPPORTED")
}
