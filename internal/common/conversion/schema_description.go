package conversion

import (
	"reflect"
	"slices"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func DataSourceSchemaFromResource(rs schema.Schema, requiredFields ...string) dsschema.Schema {
	attrs := make(map[string]dsschema.Attribute, len(rs.Attributes))
	for k, v := range rs.Attributes {
		computed := true
		required := false
		if slices.Contains(requiredFields, k) {
			computed = false
			required = true
		}
		attrs[k] = dsschema.StringAttribute{
			MarkdownDescription: v.GetMarkdownDescription(),
			Computed:            computed,
			Required:            required,
		}
	}
	ds := dsschema.Schema{
		Attributes: attrs,
	}
	UpdateSchemaDescription(&ds)
	return ds
}

func UpdateSchemaDescription[T schema.Schema | dsschema.Schema](s *T) {
	UpdateAttr(s)
}

// UpdateAttr is exported for testing purposes only and should not be used directly.
func UpdateAttr(attr any) {
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
	if strDescr != "" && strMDDescr != "" {
		panic("both descriptions exist, please fix caller: " + strDescr)
	}
	if strDescr == "" {
		fDescr.SetString(fMDDescr.String())
	} else {
		fMDDescr.SetString(fDescr.String())
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
		UpdateAttr(newPtr.Interface())
		f.SetMapIndex(k, newPtr.Elem())
	}
}

func updateNested(v reflect.Value, nestedName string) {
	f := v.FieldByName(nestedName)
	if !f.IsValid() {
		return
	}
	ptr := f.Addr()
	UpdateAttr(ptr.Interface())
}
