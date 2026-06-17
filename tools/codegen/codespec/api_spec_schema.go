package codespec

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const discriminatorExtensionKey = "x-xgen-discriminator"

const arraySemanticExtensionKey = "x-xgen-array-semantic"

const (
	arraySemanticList = "list"
	arraySemanticSet  = "set"
)

// ErrInvalidArraySemantic is returned when the x-xgen-array-semantic extension carries a value
// other than "list" or "set". It is a sentinel so callers can distinguish a malformed extension
// (which must fail generation) from a response/parameter schema that simply could not be mapped.
var ErrInvalidArraySemantic = errors.New("invalid " + arraySemanticExtensionKey + " value")

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

// GetXGenArraySemantic returns the declared array semantic ("list" or "set"), or nil if the
// extension is absent. It returns an error wrapping ErrInvalidArraySemantic when the value is
// malformed or not one of the allowed values, so generation fails wherever the property appears.
func (s *APISpecSchema) GetXGenArraySemantic() (*string, error) {
	if s.Schema.Extensions == nil {
		return nil, nil
	}

	node, ok := s.Schema.Extensions.Get(arraySemanticExtensionKey)
	if !ok || node == nil {
		return nil, nil
	}

	var value string
	if err := node.Decode(&value); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidArraySemantic, err)
	}

	if value != arraySemanticList && value != arraySemanticSet {
		return nil, fmt.Errorf("%w: %q, expected %q or %q", ErrInvalidArraySemantic, value, arraySemanticList, arraySemanticSet)
	}

	return &value, nil
}
