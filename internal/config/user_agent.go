package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

type UserAgentKey struct{} // used as contextkey

const (
	toolName = "terraform-provider-mongodbatlas"
)

var (
	userAgentProviderVersion = fmt.Sprintf("%s/%s", toolName, version.ProviderVersion)
)

func TerraformVersionUserAgentInfo(tfVersion string) string {
	return fmt.Sprintf("Terraform/%s", tfVersion)
}

func AppendToUserAgentInCtx(ctx context.Context, additionalInfo string) context.Context {
	existingUA, _ := ctx.Value(UserAgentKey{}).(string)

	if !strings.Contains(existingUA, additionalInfo) {
		if existingUA != "" {
			existingUA += " "
		}
		existingUA += additionalInfo
	}

	return context.WithValue(ctx, UserAgentKey{}, existingUA)
}

// func AppendToUserAgent(existingUA, additionalInfo string) string {
// 	// existingUA, _ := ctx.Value(UserAgentKey{}).(string)

// 	if !strings.Contains(existingUA, additionalInfo) {
// 		if existingUA != "" {
// 			existingUA += " "
// 		}
// 		existingUA += additionalInfo
// 	}
// 	return existingUA
// 	// return context.WithValue(ctx, UserAgentKey{}, existingUA)
// }
