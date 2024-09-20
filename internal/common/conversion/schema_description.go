package conversion

import (
	"reflect"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UpdateSchemaDescription(s *schema.Schema) {
	updateAttr(s)
}

func UpdateDSSchemaDescription(s *dsschema.Schema) {
	updateAttr(s)
}

func updateAttr(attr any) {
	ptr := reflect.ValueOf(attr)
	if ptr.Kind() != reflect.Ptr {
		panic("not ptr, please fix caller")
	}
	v := ptr.Elem()
	if v.Kind() != reflect.Struct {
		panic("not struct, please fix caller")
	}
	updateDesc(v)
	updateMap(v, "Attributes")
	updateMap(v, "Blocks")
	updateNested(v, "NestedObject")
}

func updateDesc(v reflect.Value) {
	fDescr, fMDDescr := v.FieldByName("Description"), v.FieldByName("MarkdownDescription")
	if !fDescr.IsValid() || !fMDDescr.IsValid() {
		return
	}
	if !fDescr.CanSet() || fDescr.Kind() != reflect.String ||
		!fMDDescr.CanSet() || fMDDescr.Kind() != reflect.String {
		panic("invalid desc fields, please fix caller")
	}
	strDescr, strMDDescr := fDescr.String(), fMDDescr.String()
	if strDescr == "" && strMDDescr != "" {
		fDescr.SetString(fMDDescr.String())
		return
	}
	if strMDDescr == "" && strDescr != "" {
		fMDDescr.SetString(fDescr.String())
		return
	}
	if strDescr != "" && strDescr != strMDDescr {
		panic("conflicting descriptions, please fix caller: " + strDescr)
	}
}

func updateMap(v reflect.Value, mapName string) {
	f := v.FieldByName(mapName)
	if !f.IsValid() {
		return
	}
	if f.Kind() != reflect.Map {
		panic("not map, please fix caller: " + mapName)
	}
	for _, k := range f.MapKeys() {
		v := f.MapIndex(k).Elem()
		newPtr := reflect.New(v.Type())
		newPtr.Elem().Set(v)
		updateAttr(newPtr.Interface())
		f.SetMapIndex(k, newPtr.Elem())
	}
}

func updateNested(v reflect.Value, nestedName string) {
	f := v.FieldByName(nestedName)
	if !f.IsValid() {
		return
	}
	ptr := f.Addr()
	updateAttr(ptr.Interface())
}
