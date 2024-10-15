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
