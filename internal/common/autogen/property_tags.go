package autogen

import (
	"reflect"
	"slices"
	"strings"
)

const (
	tagKey                    = "autogen"
	tagSensitive              = "sensitive"
	tagValOmitJSON            = "omitjson"
	tagValOmitJSONUpdate      = "omitjsonupdate"
	tagValIncludeNullOnUpdate = "includenullonupdate"
	tagValListAsMap           = "listasmap"
	tagAPIName                = "apiname" // e.g., apiname:"groupId" means JSON field is "groupId", used if the API name is different from the uncapitalized model name
)

type PropertyTags struct {
	APIName             *string
	Sensitive           bool
	OmitJSON            bool
	OmitJSONUpdate      bool
	IncludeNullOnUpdate bool
	ListAsMap           bool
}

func GetPropertyTags(field *reflect.StructField) PropertyTags {
	tags := strings.Split(field.Tag.Get(tagKey), ",")
	result := PropertyTags{
		Sensitive:           slices.Contains(tags, tagSensitive),
		OmitJSON:            slices.Contains(tags, tagValOmitJSON),
		OmitJSONUpdate:      slices.Contains(tags, tagValOmitJSONUpdate),
		IncludeNullOnUpdate: slices.Contains(tags, tagValIncludeNullOnUpdate),
		ListAsMap:           slices.Contains(tags, tagValListAsMap),
	}
	if apiName := field.Tag.Get(tagAPIName); apiName != "" {
		result.APIName = &apiName
	}
	return result
}
