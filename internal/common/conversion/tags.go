package conversion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func NewResourceTags(ctx context.Context, tags types.Map) *[]admin.ResourceTag {
	if tags.IsNull() || len(tags.Elements()) == 0 {
		return &[]admin.ResourceTag{}
	}
	elements := make(map[string]types.String, len(tags.Elements()))
	_ = tags.ElementsAs(ctx, &elements, false)
	var tagsAdmin []admin.ResourceTag
	for key, tagValue := range elements {
		tagsAdmin = append(tagsAdmin, admin.ResourceTag{
			Key:   key,
			Value: tagValue.ValueString(),
		})
	}
	return &tagsAdmin
}

func NewTFTags(tags []admin.ResourceTag) types.Map {
	if len(tags) == 0 {
		return types.MapNull(types.StringType)
	}
	typesTags := make(map[string]attr.Value, len(tags))
	for _, tag := range tags {
		typesTags[tag.Key] = types.StringValue(tag.Value)
	}
	return types.MapValueMust(types.StringType, typesTags)
}
