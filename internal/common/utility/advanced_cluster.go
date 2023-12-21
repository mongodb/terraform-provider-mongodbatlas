package utility

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spf13/cast"
)

func FormatMongoDBMajorVersion(val any) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}

	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
}

func IsListPresent(v types.List) bool {
	return !v.IsNull() && len(v.Elements()) > 0
}

func IsSetPresent(v types.Set) bool {
	return !v.IsNull() && len(v.Elements()) > 0
}
