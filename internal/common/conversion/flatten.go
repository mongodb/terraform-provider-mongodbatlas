package conversion

import "go.mongodb.org/atlas-sdk/v20231115005/admin"

func FlattenLinks(links []admin.Link) []map[string]any {
	ret := make([]map[string]any, len(links))
	for i, link := range links {
		ret[i] = map[string]any{
			"href": link.GetHref(),
			"rel":  link.GetRel(),
		}
	}
	return ret
}

func FlattenTags(tags []admin.ResourceTag) []map[string]any {
	ret := make([]map[string]any, len(tags))
	for i, tag := range tags {
		ret[i] = map[string]any{
			"key":   tag.GetKey(),
			"value": tag.GetValue(),
		}
	}
	return ret
}
