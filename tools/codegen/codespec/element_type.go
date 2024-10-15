package codespec

import "fmt"

type ElemType int

const (
	Bool ElemType = iota
	Float64
	Int64
	Number
	String
	Unknown
)

// type ElemType struct {
// 	Bool    *BoolType    `json:"bool,omitempty"`
// 	Float64 *Float64Type `json:"float64,omitempty"`
// 	Int64   *Int64Type   `json:"int64,omitempty"`
// 	List    *ListType    `json:"list,omitempty"`
// 	Number  *NumberType  `json:"number,omitempty"`
// 	String  *StringType  `json:"string,omitempty"`
// }

// type BoolType struct{}
// type Float64Type struct{}
// type Int64Type struct{}
// type StringType struct{}

// type NumberType struct{}
// type ListType struct {
// 	ElementType ElemType `json:"element_type"`
// }
// type SetType struct {
// 	ElementType ElemType `json:"element_type"`
// }

// // type ObjectType struct{} // TODO
// type MapType struct {
// 	ElementType ElemType `json:"element_type"`
// }

func (s *APISpecSchema) buildElementType() (ElemType, error) {
	switch s.Type {
	case OASTypeString:
		return String, nil
	case OASTypeBoolean:
		return Bool, nil
	case OASTypeInteger:
		return Int64, nil
	case OASTypeNumber:
		return Number, nil
	case OASTypeArray, OASTypeObject:
		return Unknown, nil // ignoring because complex element types unsupported
	default:
		return Unknown, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *APISpecSchema) buildArrayElementType() (ElemType, error) {
	if !s.Schema.Items.IsA() {
		return Unknown, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return Unknown, err
	}

	elemType, err := itemSchema.buildElementType()
	if err != nil {
		return Unknown, err
	}

	return elemType, nil
}

// func (s *APISpecSchema) buildStringElementType() (ElemType, error) {
// 	return ElemType{
// 		String: &StringType{},
// 	}, nil
// }

// func (s *APISpecSchema) buildIntegerElementType() (ElemType, error) {
// 	return ElemType{
// 		Int64: &Int64Type{},
// 	}, nil
// }

// func (s *APISpecSchema) buildBoolElementType() (ElemType, error) {
// 	return ElemType{
// 		Bool: &BoolType{},
// 	}, nil
// }

// func (s *APISpecSchema) buildNumberElementType() (ElemType, error) {
// 	if s.Schema.Format == OASFormatDouble || s.Schema.Format == OASFormatFloat {
// 		return ElemType{
// 			Float64: &Float64Type{},
// 		}, nil
// 	}

// 	return ElemType{
// 		Number: &NumberType{},
// 	}, nil
// }
