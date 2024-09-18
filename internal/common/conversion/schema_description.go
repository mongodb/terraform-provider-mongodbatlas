package conversion

import (
	"fmt"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UpdateSchemaDescription(s *schema.Schema) {
	for i := range s.Attributes {
		s.Attributes[i] = updateAttribute(s.Attributes[i])
	}
	for i := range s.Blocks {
		s.Blocks[i] = updateBlock(s.Blocks[i])
	}
}

func UpdateDSSchemaDescription(s *dsschema.Schema) {
	for i := range s.Attributes {
		s.Attributes[i] = updateAttribute(s.Attributes[i])
	}
	for i := range s.Blocks {
		s.Blocks[i] = updateBlock(s.Blocks[i])
	}
}

func updateAttribute(attr schema.Attribute) schema.Attribute {
	switch v := attr.(type) {
	case schema.StringAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case schema.BoolAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case schema.Int64Attribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v

	case dsschema.StringAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case dsschema.BoolAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case dsschema.Int64Attribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case dsschema.MapAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v
	case dsschema.ListAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		return v

	case schema.SingleNestedAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		for i := range v.Attributes {
			v.Attributes[i] = updateAttribute(v.Attributes[i])
		}
		return v
	case schema.ListNestedAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		for i := range v.NestedObject.Attributes {
			v.NestedObject.Attributes[i] = updateAttribute(v.NestedObject.Attributes[i])
		}
		return v

	case dsschema.SingleNestedAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		for i := range v.Attributes {
			v.Attributes[i] = updateAttribute(v.Attributes[i])
		}
		return v
	case dsschema.ListNestedAttribute:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		} else if v.MarkdownDescription == "" {
			v.MarkdownDescription = v.Description
		}
		for i := range v.NestedObject.Attributes {
			v.NestedObject.Attributes[i] = updateAttribute(v.NestedObject.Attributes[i])
		}
		return v
	}
	// add more attribute types as needed
	panic(fmt.Sprintf("unsupported attribute updating description, type: %T, value: %+v", attr, attr))
}

func updateBlock(block schema.Block) schema.Block {
	switch v := block.(type) {
	case schema.ListNestedBlock:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		}
		for i := range v.NestedObject.Attributes {
			v.NestedObject.Attributes[i] = updateAttribute(v.NestedObject.Attributes[i])
		}
		for i := range v.NestedObject.Blocks {
			v.NestedObject.Blocks[i] = updateBlock(v.NestedObject.Blocks[i])
		}
		return v

	case dsschema.ListNestedBlock:
		if v.Description == "" {
			v.Description = v.MarkdownDescription
		}
		for i := range v.NestedObject.Attributes {
			v.NestedObject.Attributes[i] = updateAttribute(v.NestedObject.Attributes[i])
		}
		for i := range v.NestedObject.Blocks {
			v.NestedObject.Blocks[i] = updateBlock(v.NestedObject.Blocks[i])
		}
		return v
	}
	// add more block types as needed
	panic(fmt.Sprintf("unsupported block updating description, type: %T, value: %+v", block, block))
}
