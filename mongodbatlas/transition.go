package mongodbatlas

import (
	"hash/crc32"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func valRegion(reg any, opt ...string) (string, error) {
	return config.ValRegion(reg, opt...)
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	return config.RemoveLabel(list, item)
}

func pointer[T any](x T) *T {
	return &x
}

func intPtr(v int) *int {
	if v != 0 {
		return &v
	}
	return nil
}

func stringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

// HashCodeString hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func HashCodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func expandStringList(list []any) (res []string) {
	return config.ExpandStringList(list)
}
