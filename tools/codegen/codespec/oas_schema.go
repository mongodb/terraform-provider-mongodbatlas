package codespec

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var errInvalidSchema = fmt.Errorf("invalid schema")

type OASSchema struct {
	Schema *base.Schema
	Type   string
	Format string
}

type OASResource struct {
	Description      *string
	CreateOp         *high.Operation
	ReadOp           *high.Operation
	UpdateOp         *high.Operation
	DeleteOp         *high.Operation
	CommonParameters []*high.Parameter
}

func (s *OASSchema) GetComputability(name string) ComputedOptionalRequired {
	for _, prop := range s.Schema.Required {
		if name == prop {
			return Required
		}
	}

	return ComputedOptional
}

func (s *OASSchema) GetDeprecationMessage() *string {
	if s.Schema.Deprecated == nil || !(*s.Schema.Deprecated) {
		return nil
	}

	deprecationMessage := "This attribute has been deprecated"

	return &deprecationMessage
}

func (s *OASSchema) GetDescription() *string {
	if s.Schema.Description == "" {
		return nil
	}

	return &s.Schema.Description
}

func (s *OASSchema) IsSensitive() *bool {
	isSensitive := s.Format == OASFormatPassword

	if !isSensitive {
		return nil
	}

	return &isSensitive
}
