package config

import (
	"context"
	"fmt"
	"log"
	"maps"
	"sort"
	"strings"
)

const (
	UserAgentKeyName                    = "Name"
	UserAgentKeyOperation               = "Operation"
	UserAgentKeyModuleName              = "ModuleName"
	UserAgentKeyModuleVersion           = "ModuleVersion"
	UserAgentOperationValueCreate       = "create"
	UserAgentOperationValueRead         = "read"
	UserAgentOperationValueUpdate       = "update"
	UserAgentOperationValueDelete       = "delete"
	UserAgentOperationValueImport       = "import"
	UserAgentOperationValuePlanModify   = "plan-modify"
	UserAgentOperationValueUpgradeState = "upgrade-state"
	UserAgentOperationValueMoveState    = "move-state"
)

// UserAgentExtra holds additional metadata to be appended to the User-Agent header and context.
type UserAgentExtra struct {
	Extras        map[string]string
	Name          string
	Operation     string
	ModuleName    string
	ModuleVersion string
}

func userAgentNameValue(name string) string {
	return strings.TrimPrefix(name, "mongodbatlas_")
}

// Combine returns a new UserAgentExtra by merging the receiver with another.
// Non-empty fields in 'other' take precedence over the receiver's fields.
func (e UserAgentExtra) Combine(other UserAgentExtra) UserAgentExtra {
	name := e.Name
	if other.Name != "" {
		name = other.Name
	}
	operation := e.Operation
	if other.Operation != "" {
		operation = other.Operation
	}
	moduleName := e.ModuleName
	if other.ModuleName != "" {
		moduleName = other.ModuleName
	}
	moduleVersion := e.ModuleVersion
	if other.ModuleVersion != "" {
		moduleVersion = other.ModuleVersion
	}
	var newExtras map[string]string
	if e.Extras != nil {
		newExtras = map[string]string{}
		maps.Copy(newExtras, e.Extras)
	}
	if other.Extras != nil {
		if newExtras == nil {
			newExtras = map[string]string{}
		}
		maps.Copy(newExtras, other.Extras)
	}
	return UserAgentExtra{
		Name:          name,
		Operation:     operation,
		ModuleName:    moduleName,
		ModuleVersion: moduleVersion,
		Extras:        newExtras,
	}
}

// ToHeaderValue returns a string representation suitable for use as a User-Agent header value.
// If oldHeader is non-empty, it is prepended to the new value.
func (e UserAgentExtra) ToHeaderValue(ctx context.Context, oldHeader string) string {
	parts := map[string]string{}
	addPart := func(key, part string) {
		if part == "" {
			return
		}
		if existing, found := parts[key]; found {
			log.Printf("[WARN] Replaced UserAgent key %s: %s -> %s", key, existing, part)
		}
		parts[key] = part
	}
	// Start with Extras to avoid malicious usage of known keys
	for k, v := range e.Extras {
		addPart(k, v)
	}
	addPart(UserAgentKeyName, e.Name)
	addPart(UserAgentKeyOperation, e.Operation)
	addPart(UserAgentKeyModuleName, e.ModuleName)
	addPart(UserAgentKeyModuleVersion, e.ModuleVersion)
	partsLen := len(parts)
	if partsLen == 0 {
		return oldHeader
	}
	sortedKeys := make([]string, 0, partsLen)
	for k := range parts {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	newStrings := make([]string, partsLen)
	for i, k := range sortedKeys {
		newStrings[i] = fmt.Sprintf("%s/%s", k, parts[k])
	}
	newPart := strings.Join(newStrings, " ")
	if oldHeader == "" {
		return newPart
	}
	return fmt.Sprintf("%s %s", oldHeader, newPart)
}

type UserAgentKey string

const (
	UserAgentExtraKey = UserAgentKey("user-agent-extra")
	UserAgentHeader   = "User-Agent"
)

// ReadUserAgentExtra retrieves the UserAgentExtra from the context if present.
// Returns a pointer to the UserAgentExtra, or nil if not set or of the wrong type.
// Logs a warning if the value is not of the expected type.
func ReadUserAgentExtra(ctx context.Context) *UserAgentExtra {
	extra := ctx.Value(UserAgentExtraKey)
	if extra == nil {
		return nil
	}
	if userAgentExtra, ok := extra.(UserAgentExtra); ok {
		return &userAgentExtra
	}
	log.Printf("[WARN] UserAgentExtra in context is not of type UserAgentExtra, got %v", extra)
	return nil
}

// AddUserAgentExtra returns a new context with UserAgentExtra merged into any existing value.
// If a UserAgentExtra is already present in the context, the fields of 'extra' will override non-empty fields.
func AddUserAgentExtra(ctx context.Context, extra UserAgentExtra) context.Context {
	oldExtra := ReadUserAgentExtra(ctx)
	if oldExtra == nil {
		return context.WithValue(ctx, UserAgentExtraKey, extra)
	}
	newExtra := oldExtra.Combine(extra)
	return context.WithValue(ctx, UserAgentExtraKey, newExtra)
}
