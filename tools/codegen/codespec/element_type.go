package codespec

import "fmt"

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
		return String, nil // complex element types are unsupported so this defaults to string for now to provide best effort generation
	default:
		return Unknown, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
