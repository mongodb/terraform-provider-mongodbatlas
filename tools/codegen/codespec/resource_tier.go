package codespec

import "fmt"

type ResourceTier string

const (
	ResourceTierProd     ResourceTier = "prod"
	ResourceTierInternal ResourceTier = "internal"
)

func ParseResourceTier(value string) (*ResourceTier, error) {
	if value == "" {
		return nil, nil
	}
	tier := ResourceTier(value)
	switch tier {
	case ResourceTierProd, ResourceTierInternal:
		return &tier, nil
	default:
		return nil, fmt.Errorf("expected %q or %q, got %q", ResourceTierProd, ResourceTierInternal, value)
	}
}
