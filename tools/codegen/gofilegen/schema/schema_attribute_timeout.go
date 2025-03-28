package schema

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
)

type timeoutAttributeGenerator struct {
	timeouts codespec.TimeoutsAttribute
}

func (s *timeoutAttributeGenerator) AttributeCode() CodeStatement {
	var optionProperties string
	for _, op := range s.timeouts.ConfigurableTimeouts {
		switch op {
		case codespec.Create:
			optionProperties += "Create: true,\n"
		case codespec.Update:
			optionProperties += "Update: true,\n"
		case codespec.Delete:
			optionProperties += "Delete: true,\n"
		case codespec.Read:
			optionProperties += "Read: true,\n"
		}
	}
	return CodeStatement{
		Code: fmt.Sprintf(`"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
			%s
		})`, optionProperties),
		Imports: []string{"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"},
	}
}
