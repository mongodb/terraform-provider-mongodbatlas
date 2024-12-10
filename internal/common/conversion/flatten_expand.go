package conversion

import (
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	adminPreview "github.com/mongodb/atlas-sdk-go/admin" // TODO: replace usage with latest once cipher config changes are in prod
)

func FlattenLinks(links []admin.Link) []map[string]string {
	ret := make([]map[string]string, len(links))
	for i, link := range links {
		ret[i] = map[string]string{
			"href": link.GetHref(),
			"rel":  link.GetRel(),
		}
	}
	return ret
}

func FlattenTags(tags []admin.ResourceTag) []map[string]string {
	ret := make([]map[string]string, len(tags))
	for i, tag := range tags {
		ret[i] = map[string]string{
			"key":   tag.GetKey(),
			"value": tag.GetValue(),
		}
	}
	return ret
}

// TODO: remove once new SDK is available with cipher config changes
func FlattenTagsPreview(tags []adminPreview.ResourceTag) []map[string]string {
	ret := make([]map[string]string, len(tags))
	for i, tag := range tags {
		ret[i] = map[string]string{
			"key":   tag.GetKey(),
			"value": tag.GetValue(),
		}
	}
	return ret
}

func ExpandTagsFromSetSchemaV220240805(d *schema.ResourceData) *[]admin20240805.ResourceTag {
	list := d.Get("tags").(*schema.Set)
	ret := make([]admin20240805.ResourceTag, list.Len())
	for i, item := range list.List() {
		tag := item.(map[string]any)
		ret[i] = admin20240805.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
	}
	return &ret
}

func ExpandTagsFromSetSchema(d *schema.ResourceData) *[]admin.ResourceTag {
	list := d.Get("tags").(*schema.Set)
	ret := make([]admin.ResourceTag, list.Len())
	for i, item := range list.List() {
		tag := item.(map[string]any)
		ret[i] = admin.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
	}
	return &ret
}

func ExpandStringList(list []any) (res []string) {
	for _, v := range list {
		res = append(res, v.(string))
	}
	return
}

func ExpandStringListFromSetSchema(set *schema.Set) []string {
	res := make([]string, set.Len())
	for i, v := range set.List() {
		res[i] = v.(string)
	}
	return res
}
