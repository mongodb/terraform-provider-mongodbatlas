package thirdpartyintegration

import (
	"context"
	"slices"

	"go.mongodb.org/atlas-sdk/v20250312028/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/sdkv2config"
)

// writeOnlySecret is a secret attribute with a write-only alternative, <attr>_wo, never stored in plan or state.
// Terraform can't diff <attr>_wo, so <attr>_wo_version tells when to send a new value.
type writeOnlySecret struct {
	set  func(integration *admin.ThirdPartyIntegration, value *string)
	attr string
	// alwaysSend resends <attr>_wo on every update. updateIntegrationFromSchema never reuses the api_key
	// returned by Atlas, so every update must carry the configured value.
	alwaysSend bool
}

var writeOnlySecrets = []writeOnlySecret{
	{
		attr:       "api_key",
		alwaysSend: true,
		set:        func(i *admin.ThirdPartyIntegration, v *string) { i.ApiKey = v },
	},
	{
		attr: "service_key",
		set:  func(i *admin.ThirdPartyIntegration, v *string) { i.ServiceKey = v },
	},
	{
		attr: "routing_key",
		set:  func(i *admin.ThirdPartyIntegration, v *string) { i.RoutingKey = v },
	},
	{
		attr: "secret",
		set:  func(i *admin.ThirdPartyIntegration, v *string) { i.Secret = v },
	},
	{
		attr: "microsoft_teams_webhook_url",
		set:  func(i *admin.ThirdPartyIntegration, v *string) { i.MicrosoftTeamsWebhookUrl = v },
	},
	{
		attr: "password",
		set:  func(i *admin.ThirdPartyIntegration, v *string) { i.Password = v },
	},
}

func (s writeOnlySecret) woAttr() string        { return s.attr + "_wo" }
func (s writeOnlySecret) woVersionAttr() string { return s.attr + "_wo_version" }

func isWriteOnlySecret(attr string) bool {
	return slices.ContainsFunc(writeOnlySecrets, func(s writeOnlySecret) bool { return s.attr == attr })
}

// addWriteOnlySecretAttributes adds <attr>_wo and <attr>_wo_version for each secret.
func addWriteOnlySecretAttributes(attrs map[string]*schema.Schema) {
	for _, s := range writeOnlySecrets {
		attrs[s.attr].ConflictsWith = []string{s.woAttr()}
		attrs[s.woAttr()] = &schema.Schema{
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			WriteOnly:     true,
			ConflictsWith: []string{s.attr},
			RequiredWith:  []string{s.woVersionAttr()},
		}
		attrs[s.woVersionAttr()] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			RequiredWith: []string{s.woAttr()},
			ValidateFunc: validation.IntAtLeast(1),
		}
	}
}

// writeOnlyValues returns the configured <attr>_wo values keyed by attr.
// Write-only values are null in plan and state, so they are read from raw config.
func writeOnlyValues(d *schema.ResourceData) map[string]string {
	values := make(map[string]string)
	for _, s := range writeOnlySecrets {
		if v := sdkv2config.String(d, s.woAttr()); v.Set() && v.ValueString() != "" {
			values[s.attr] = v.ValueString()
		}
	}
	return values
}

// setWriteOnlyValues sets each write-only value whose version changed, which is every configured value on create.
func setWriteOnlyValues(d *schema.ResourceData, values map[string]string, integration *admin.ThirdPartyIntegration) {
	for _, s := range writeOnlySecrets {
		if v, ok := values[s.attr]; ok && (s.alwaysSend || d.HasChange(s.woVersionAttr())) {
			s.set(integration, &v)
		}
	}
}

// clearSecretsWithWriteOnly blanks each secret set through <attr>_wo, so state keeps neither the redacted
// value returned by Atlas nor a plaintext value from before switching to <attr>_wo.
func clearSecretsWithWriteOnly(d *schema.ResourceData, integrationMap map[string]any) {
	for _, s := range writeOnlySecrets {
		if _, ok := d.GetOk(s.woVersionAttr()); ok {
			integrationMap[s.attr] = ""
		}
	}
}

// planClearSecretsWithWriteOnly marks a plaintext secret left in state as changing when switching to <attr>_wo,
// so the plan matches the blanking done by clearSecretsWithWriteOnly. SDKv2 can't plan "" for a Computed
// attribute, so it is planned as unknown.
func planClearSecretsWithWriteOnly(_ context.Context, d *schema.ResourceDiff, _ any) error {
	for _, s := range writeOnlySecrets {
		if _, ok := d.GetOk(s.woVersionAttr()); !ok {
			continue
		}
		// Optional-only secrets already plan "" once removed from config, and SetNewComputed only accepts Computed ones.
		if !d.NewValueKnown(s.attr) || d.Get(s.attr).(string) == "" {
			continue
		}
		if err := d.SetNewComputed(s.attr); err != nil {
			return err
		}
	}
	return nil
}
