package codespec

import (
	"log"
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const discriminatorExtensionKey = "x-xgen-discriminator"

// DiscriminatorExtension represents the raw x-xgen-discriminator extension as declared in the OpenAPI spec.
type DiscriminatorExtension struct {
	Mapping      map[string]DiscriminatorExtensionType `yaml:"mapping"`
	PropertyName string                                `yaml:"propertyName"`
}

type DiscriminatorExtensionType struct {
	Properties []string `yaml:"properties,omitempty"`
	Required   []string `yaml:"required,omitempty"`
}

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
}

func (s *APISpecSchema) GetComputability(name string) ComputedOptionalRequired {
	if slices.Contains(s.Schema.Required, name) {
		return Required
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

func (s *APISpecSchema) IsSensitive() bool {
	return s.Schema.Format == OASFormatPassword
}

// GetXGenDiscriminator extracts the raw x-xgen-discriminator extension from the schema.
// Returns nil if the extension is absent or cannot be decoded.
func (s *APISpecSchema) GetXGenDiscriminator() *DiscriminatorExtension {
	if s.Schema.Extensions == nil {
		return nil
	}

	node, ok := s.Schema.Extensions.Get(discriminatorExtensionKey)
	if !ok || node == nil {
		return nil
	}

	var result DiscriminatorExtension
	if err := node.Decode(&result); err != nil {
		log.Printf("[WARN] Failed to decode %s extension: %s", discriminatorExtensionKey, err)
		return nil
	}

	return &result
}
