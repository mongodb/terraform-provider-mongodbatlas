package conversion

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type TPFSchema interface {
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
}

type TPFSrc interface {
	GetAttribute(context.Context, path.Path, any) diag.Diagnostics
}

// AttributePathValue retrieves the value for src (state/plan/config) @ attributePath with converted path.Path, schema is needed to get the correct types.XXX (String/Object/etc.)
func AttributePathValue(ctx context.Context, diags *diag.Diagnostics, attributePath *tftypes.AttributePath, src TPFSrc, schema TPFSchema) (attr.Value, path.Path) {
	convertedPath, localDiags := AttributePath(ctx, attributePath, schema)
	diags.Append(localDiags...)
	if diags.HasError() {
		return nil, convertedPath
	}
	attrType, err := schema.TypeAtTerraformPath(ctx, attributePath)
	if err != nil {
		diags.AddError("Unable to get type for attribute path", fmt.Sprintf("%s: %s", attributePath.String(), err))
		return nil, convertedPath
	}
	attrValue := attrType.ValueType(ctx)
	if localDiags := src.GetAttribute(ctx, convertedPath, &attrValue); localDiags.HasError() {
		diags.Append(localDiags...)
		return nil, convertedPath
	}
	return attrValue, convertedPath
}

const keyValue = "ElementKeyValue("

var prefixes = map[string]func(path.Path, string) (path.Path, error){
	"AttributeName(": func(p path.Path, s string) (path.Path, error) {
		return p.AtName(s), nil
	},
	"ElementKeyString(": func(p path.Path, s string) (path.Path, error) {
		return p.AtMapKey(s), nil
	},
	"ElementKeyInt(": func(p path.Path, s string) (path.Path, error) {
		number, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Sprintf("could not convert %s to int: %v", s, err))
		}
		return p.AtListIndex(number), nil
	},
	keyValue: func(p path.Path, s string) (path.Path, error) {
		if strings.HasPrefix(s, "tftypes.String<") {
			s = strings.TrimPrefix(s, `tftypes.String<"`)
			s = strings.TrimSuffix(s, `">`)
		} else {
			return path.Empty(), fmt.Errorf("could not convert %s at path %s", s, p.String())
		}
		return p.AtSetValue(types.StringValue(s)), nil
	},
}

func ConvertAttributePath(in tftypes.AttributePath) (path.Path, error) {
	tpfPath := path.Empty()
	inString := in.String()
	parts := strings.Split(inString, ".")
	var err error
	var done bool
	for i, part := range parts {
		tpfPath, done, err = addStep(part, strings.Join(parts[i:], "."), tpfPath)
		if err != nil {
			return path.Empty(), err
		}
		if done {
			return tpfPath, nil
		}
	}
	return tpfPath, nil
}

func addStep(part, remaingParts string, tpfPath path.Path) (path.Path, bool, error) {
	var err error
	for prefix, replacer := range prefixes {
		if !strings.HasPrefix(part, prefix) {
			continue
		}
		var done bool
		if prefix == keyValue {
			part = remaingParts
			done = true
		}
		part = strings.TrimPrefix(part, prefix)
		part = strings.TrimSuffix(part, ")")
		part = strings.Trim(part, `"`)
		tpfPath, err = replacer(tpfPath, part)
		if err != nil {
			return path.Empty(), false, fmt.Errorf("could not convert %s: %w", part, err)
		}
		return tpfPath, done, nil
	}
	return path.Empty(), false, fmt.Errorf("unknown prefix %s to convert, current path %s", part, tpfPath.String())
}

func AttributePath(ctx context.Context, tfType *tftypes.AttributePath, schema TPFSchema) (path.Path, diag.Diagnostics) {
	p, err := ConvertAttributePath(*tfType)
	if err != nil {
		panic(fmt.Sprintf("could not convert TF path %s to TPF path: %v", tfType.String(), err))
		// return path.Empty(), diag.Diagnostics{
		// 	diag.NewErrorDiagnostic(
		// 		"Unable to Convert Attribute Path to TPF Path",
		// 		fmt.Sprintf("An unexpected error occurred while trying to convert an attribute path. %s", err),
		// 	),
		// }
	}
	return p, nil
}
