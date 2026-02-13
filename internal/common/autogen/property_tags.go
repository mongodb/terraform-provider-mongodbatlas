package autogen

import (
	"reflect"
	"slices"
	"strings"
)

const (
	tagKey                        = "autogen"
	tagSensitive                  = "sensitive"
	tagValOmitJSON                = "omitjson"
	tagValOmitJSONUpdate          = "omitjsonupdate"
	tagValSendNullAsNullOnUpdate  = "sendnullasnullonupdate"
	tagValSendNullAsEmptyOnUpdate = "sendnullasemptyonupdate"
	tagValListAsMap               = "listasmap"
	tagValSkipStateListMerge      = "skipstatelistmerge"
	tagAPIName                    = "apiname" // e.g., apiname:"groupId" means JSON field is "groupId", used if the API name is different from the uncapitalized model name
)

type PropertyTags struct {
	APIName                 *string
	Sensitive               bool
	OmitJSON                bool
	OmitJSONUpdate          bool
	SendNullAsNullOnUpdate  bool
	SendNullAsEmptyOnUpdate bool
	ListAsMap               bool
	SkipStateListMerge      bool
}

func GetPropertyTags(field *reflect.StructField) PropertyTags {
	tags := strings.Split(field.Tag.Get(tagKey), ",")
	result := PropertyTags{
		Sensitive:               slices.Contains(tags, tagSensitive),
		OmitJSON:                slices.Contains(tags, tagValOmitJSON),
		OmitJSONUpdate:          slices.Contains(tags, tagValOmitJSONUpdate),
		SendNullAsNullOnUpdate:  slices.Contains(tags, tagValSendNullAsNullOnUpdate),
		SendNullAsEmptyOnUpdate: slices.Contains(tags, tagValSendNullAsEmptyOnUpdate),
		ListAsMap:               slices.Contains(tags, tagValListAsMap),
		SkipStateListMerge:      slices.Contains(tags, tagValSkipStateListMerge),
	}
	if apiName := field.Tag.Get(tagAPIName); apiName != "" {
		result.APIName = &apiName
	}
	return result
}
