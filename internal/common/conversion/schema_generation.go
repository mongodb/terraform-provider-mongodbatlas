package conversion

import (
	"reflect"
	"slices"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func DataSourceSchemaFromResource(rs schema.Schema, requiredFields []string, overridenFields map[string]dsschema.Attribute) dsschema.Schema {
	blocks := convertBlocks(rs.Blocks, requiredFields)
	attrs := convertAttrs(rs.Attributes, requiredFields)
	for name, attr := range overridenFields {
		if attr == nil {
			delete(attrs, name)
		} else {
			attrs[name] = attr
		}
	}
	ds := dsschema.Schema{Attributes: attrs, Blocks: blocks}
	UpdateSchemaDescription(&ds)
	return ds
}

func PluralDataSourceSchemaFromResource(rs schema.Schema, requiredFields []string) dsschema.Schema {
	blocks := convertBlocks(rs.Blocks, nil)
	if len(blocks) > 0 {
		panic("blocks not supported yet in auto-generated plural data source schema as they can't go in ListNestedAttribute")
	}
	resultAttrs := convertAttrs(rs.Attributes, nil)
	rootAttrs := convertAttrs(rs.Attributes, requiredFields)
	for name := range rootAttrs {
		if !slices.Contains(requiredFields, name) {
			delete(rootAttrs, name)
		}
	}
	rootAttrs["results"] = dsschema.ListNestedAttribute{
		Computed: true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: resultAttrs,
		},
		MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
	}
	ds := dsschema.Schema{Attributes: rootAttrs}
	UpdateSchemaDescription(&ds)
	return ds
}

func UpdateSchemaDescription[T schema.Schema | dsschema.Schema](s *T) {
	UpdateAttr(s)
}

var convertMappings = map[string]reflect.Type{
	"StringAttribute":       reflect.TypeOf(dsschema.StringAttribute{}),
	"BoolAttribute":         reflect.TypeOf(dsschema.BoolAttribute{}),
	"Int64Attribute":        reflect.TypeOf(dsschema.Int64Attribute{}),
	"Float64Attribute":      reflect.TypeOf(dsschema.Float64Attribute{}),
	"MapAttribute":          reflect.TypeOf(dsschema.MapAttribute{}),
	"SingleNestedAttribute": reflect.TypeOf(dsschema.SingleNestedAttribute{}),
	"ListNestedAttribute":   reflect.TypeOf(dsschema.ListNestedAttribute{}),
	"SetNestedAttribute":    reflect.TypeOf(dsschema.SetNestedAttribute{}),
	"ListAttribute":         reflect.TypeOf(dsschema.ListAttribute{}),
	"SetNestedBlock":        reflect.TypeOf(dsschema.SetNestedBlock{}),
	"SetAttribute":          reflect.TypeOf(dsschema.SetAttribute{}),
}

var convertNestedMappings = map[string]reflect.Type{
	"NestedAttributeObject": reflect.TypeOf(dsschema.NestedAttributeObject{}),
	"NestedBlockObject":     reflect.TypeOf(dsschema.NestedBlockObject{}),
}

func convertAttrs(rsAttrs map[string]schema.Attribute, requiredFields []string) map[string]dsschema.Attribute {
	const ignoreField = "timeouts"
	if rsAttrs == nil {
		return nil
	}
	dsAttrs := make(map[string]dsschema.Attribute, len(rsAttrs))
	for name, attr := range rsAttrs {
		if name == ignoreField {
			continue
		}
		dsAttrs[name] = convertElement(name, attr, requiredFields).(dsschema.Attribute)
	}
	return dsAttrs
}

func convertBlocks(rsBlocks map[string]schema.Block, requiredFields []string) map[string]dsschema.Block {
	if rsBlocks == nil {
		return nil
	}
	dsBlocks := make(map[string]dsschema.Block, len(rsBlocks))
	for name, block := range rsBlocks {
		dsBlocks[name] = convertElement(name, block, requiredFields).(dsschema.Block)
	}
	return dsBlocks
}

func convertElement(name string, element any, requiredFields []string) any {
	computed := true
	required := false
	if slices.Contains(requiredFields, name) {
		computed = false
		required = true
	}
	vSrc := reflect.ValueOf(element)
	tSrc := reflect.TypeOf(element)
	tDest := convertMappings[tSrc.Name()]
	if tDest == nil {
		panic("attribute type not support yet, add it to convertMappings: " + tSrc.Name())
	}
	vDest := reflect.New(tDest).Elem()
	vDest.FieldByName("MarkdownDescription").Set(vSrc.FieldByName("MarkdownDescription"))
	vDest.FieldByName("DeprecationMessage").Set(vSrc.FieldByName("DeprecationMessage"))
	if fSensitive := vDest.FieldByName("Sensitive"); fSensitive.CanSet() {
		fSensitive.Set(vSrc.FieldByName("Sensitive"))
	}
	if fComputed := vDest.FieldByName("Computed"); fComputed.CanSet() {
		fComputed.SetBool(computed)
	}
	if fRequired := vDest.FieldByName("Required"); fRequired.CanSet() {
		fRequired.SetBool(required)
	}
	if fElementType := vDest.FieldByName("ElementType"); fElementType.CanSet() {
		fElementType.Set(vSrc.FieldByName("ElementType"))
	}
	if fAttributes := vDest.FieldByName("Attributes"); fAttributes.CanSet() {
		attrsSrc := vSrc.FieldByName("Attributes").Interface().(map[string]schema.Attribute)
		fAttributes.Set(reflect.ValueOf(convertAttrs(attrsSrc, nil)))
	}
	if fNested := vDest.FieldByName("NestedObject"); fNested.CanSet() {
		tNested := convertNestedMappings[fNested.Type().Name()]
		if tNested == nil {
			panic("nested type not support yet, add it to convertNestedMappings: " + fNested.Type().Name())
		}
		attrsSrc := vSrc.FieldByName("NestedObject").FieldByName("Attributes").Interface().(map[string]schema.Attribute)
		vNested := reflect.New(tNested).Elem()
		vNested.FieldByName("Attributes").Set(reflect.ValueOf(convertAttrs(attrsSrc, nil)))
		fNested.Set(vNested)
	}
	return vDest.Interface()
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
