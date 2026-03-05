package schema

import (
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type timeoutAttributeGenerator struct {
	timeouts codespec.TimeoutsAttribute
}

func (s *timeoutAttributeGenerator) AttributeCode() (CodeStatement, error) {
	var optionProperties strings.Builder
	for _, op := range s.timeouts.ConfigurableTimeouts {
		switch op {
		case codespec.Create:
			optionProperties.WriteString("Create: true,\n")
		case codespec.Update:
			optionProperties.WriteString("Update: true,\n")
		case codespec.Delete:
			optionProperties.WriteString("Delete: true,\n")
		case codespec.Read:
			optionProperties.WriteString("Read: true,\n")
		}
	}
	return CodeStatement{
		Code: fmt.Sprintf(`"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
			%s
		})`, optionProperties.String()),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"},
	}, nil
}
