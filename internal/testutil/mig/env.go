package mig

import "os"

func VersionConstraint() string {
	return os.Getenv("MONGODB_ATLAS_LAST_VERSION")
}
