package codespec

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var errInvalidSchema = fmt.Errorf("invalid schema")

type APISpecSchema struct {
	Schema *base.Schema
	Type   string
}

type APISpecResource struct {
	Description        *string
	DeprecationMessage *string
	CreateOp           *high.Operation
	ReadOp             *high.Operation
	UpdateOp           *high.Operation
	DeleteOp           *high.Operation
	CommonParameters   []*high.Parameter
}

func (s *APISpecSchema) GetComputability(name string) ComputedOptionalRequired {
	for _, prop := range s.Schema.Required {
		if name == prop {
			return Required
		}
	}

	return Optional
}

func (s *APISpecSchema) GetDeprecationMessage() *string {
	if s.Schema.Deprecated == nil || !(*s.Schema.Deprecated) {
		return nil
	}

	deprecationMessage := "This attribute has been deprecated"

	return &deprecationMessage
}

func (s *APISpecSchema) GetDescription() *string {
	if s.Schema.Description == "" {
		return nil
	}

	return &s.Schema.Description
}

func (s *APISpecSchema) IsSensitive() *bool {
	isSensitive := s.Schema.Format == OASFormatPassword

	if !isSensitive {
		return nil
	}

	return &isSensitive
}
