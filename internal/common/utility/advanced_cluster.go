package utility

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

func FormatMongoDBMajorVersion(val any) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}

	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
}
