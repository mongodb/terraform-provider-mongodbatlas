package acc

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func PathFromString(pathStr string) tfjsonpath.Path {
	segments := strings.Split(pathStr, ".")
	var path tfjsonpath.Path

	for i, seg := range segments {
		if i == 0 { // Initial step (New())
			if index, err := strconv.Atoi(seg); err == nil {
				path = tfjsonpath.New(index)
			} else {
				path = tfjsonpath.New(seg)
			}
		} else { // Subsequent steps (AtMapKey()/AtSliceIndex())
			if index, err := strconv.Atoi(seg); err == nil {
				path = path.AtSliceIndex(index)
			} else {
				path = path.AtMapKey(seg)
			}
		}
	}
	return path
}
