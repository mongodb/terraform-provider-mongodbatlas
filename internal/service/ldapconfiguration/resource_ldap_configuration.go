package ldapconfiguration

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorCreate   = "error creating MongoDB LDAPConfiguration (%s): %s"
	errorUpdate   = "error updating MongoDB LDAPConfiguration (%s): %s"
	errorRead     = "error reading MongoDB LDAPConfiguration (%s): %s"
	errorDelete   = "error deleting MongoDB LDAPConfiguration (%s): %s"
	errorSettings = "error setting `%s` for LDAPConfiguration(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"authorization_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  636,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"authz_query_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_to_dn_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"substitution": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ldap_query": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	ldap := new(admin.LDAPSecuritySettings)

	if v, ok := d.GetOk("authentication_enabled"); ok {
		ldap.AuthenticationEnabled = conversion.Pointer(v.(bool))
	}

	if v, ok := d.GetOk("authorization_enabled"); ok {
		ldap.AuthorizationEnabled = conversion.Pointer(v.(bool))
	}

	if v, ok := d.GetOk("hostname"); ok {
		ldap.Hostname = conversion.Pointer(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		ldap.Port = conversion.Pointer(v.(int))
	}

	if v, ok := d.GetOk("bind_username"); ok {
		ldap.BindUsername = conversion.Pointer(v.(string))
	}

	if v, ok := d.GetOk("bind_password"); ok {
		ldap.BindPassword = conversion.Pointer(v.(string))
	}

	if v, ok := d.GetOk("ca_certificate"); ok {
		ldap.CaCertificate = conversion.Pointer(v.(string))
	}

	if v, ok := d.GetOk("authz_query_template"); ok {
		ldap.AuthzQueryTemplate = conversion.Pointer(v.(string))
	}

	if v, ok := d.GetOk("user_to_dn_mapping"); ok {
		ldap.UserToDNMapping = expandDNMapping(v.([]any))
	}

	params := &admin.UserSecurity{
		Ldap: ldap,
	}
	_, _, err := connV2.LDAPConfigurationApi.SaveLdapConfiguration(ctx, projectID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorCreate, projectID, err))
	}
	d.SetId(projectID)
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	resp, httpResp, err := connV2.LDAPConfigurationApi.GetLdapConfiguration(context.Background(), d.Id()).Execute()
	if err != nil {
		if validate.StatusNotFound(httpResp) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(errorRead, d.Id(), err))
	}

	if err = d.Set("authentication_enabled", resp.Ldap.GetAuthenticationEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "authentication_enabled", d.Id(), err))
	}
	if err = d.Set("authorization_enabled", resp.Ldap.GetAuthorizationEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "authorization_enabled", d.Id(), err))
	}
	if err = d.Set("hostname", resp.Ldap.GetHostname()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "hostname", d.Id(), err))
	}
	if err = d.Set("port", resp.Ldap.GetPort()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "port", d.Id(), err))
	}
	if err = d.Set("bind_username", resp.Ldap.GetBindUsername()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "bind_username", d.Id(), err))
	}
	if err = d.Set("ca_certificate", resp.Ldap.GetCaCertificate()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "ca_certificate", d.Id(), err))
	}
	if err = d.Set("authz_query_template", resp.Ldap.GetAuthzQueryTemplate()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "authz_query_template", d.Id(), err))
	}
	if err = d.Set("user_to_dn_mapping", flattenDNMapping(resp.Ldap.GetUserToDNMapping())); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "user_to_dn_mapping", d.Id(), err))
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ldap := new(admin.LDAPSecuritySettings)

	if d.HasChange("authentication_enabled") {
		ldap.AuthenticationEnabled = conversion.Pointer(d.Get("authentication_enabled").(bool))
	}

	if d.HasChange("authorization_enabled") {
		ldap.AuthorizationEnabled = conversion.Pointer(d.Get("authorization_enabled").(bool))
	}

	if d.HasChange("hostname") {
		ldap.Hostname = conversion.Pointer(d.Get("hostname").(string))
	}

	if d.HasChange("port") {
		ldap.Port = conversion.Pointer(d.Get("port").(int))
	}

	if d.HasChange("bind_username") {
		ldap.BindUsername = conversion.Pointer(d.Get("bind_username").(string))
	}

	if d.HasChange("bind_password") {
		ldap.BindPassword = conversion.Pointer(d.Get("bind_password").(string))
	}

	if d.HasChange("ca_certificate") {
		ldap.CaCertificate = conversion.Pointer(d.Get("ca_certificate").(string))
	}

	if d.HasChange("authz_query_template") {
		ldap.AuthzQueryTemplate = conversion.Pointer(d.Get("authz_query_template").(string))
	}

	if d.HasChange("user_to_dn_mapping") {
		ldap.UserToDNMapping = expandDNMapping(d.Get("user_to_dn_mapping").([]any))
	}

	params := &admin.UserSecurity{
		Ldap: ldap,
	}
	_, _, err := connV2.LDAPConfigurationApi.SaveLdapConfiguration(ctx, d.Id(), params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorUpdate, d.Id(), err))
	}
	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	params := &admin.UserSecurity{
		Ldap: &admin.LDAPSecuritySettings{
			AuthenticationEnabled: conversion.Pointer(false),
			AuthorizationEnabled:  conversion.Pointer(false),
		},
	}
	_, _, err := connV2.LDAPConfigurationApi.SaveLdapConfiguration(ctx, d.Id(), params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorDelete, d.Id(), err))
	}
	return nil
}

func expandDNMapping(p []any) *[]admin.UserToDNMapping {
	mappings := make([]admin.UserToDNMapping, len(p))
	for k, v := range p {
		mapping := v.(map[string]any)
		mappings[k] = admin.UserToDNMapping{
			Match:        mapping["match"].(string),
			Substitution: conversion.StringPtr(mapping["substitution"].(string)),
			LdapQuery:    conversion.StringPtr(mapping["ldap_query"].(string)),
		}
	}
	return &mappings
}

func flattenDNMapping(mappings []admin.UserToDNMapping) []map[string]string {
	ret := make([]map[string]string, len(mappings))
	for i := range mappings {
		mapping := &mappings[i]
		ret[i] = map[string]string{
			"match":        mapping.GetMatch(),
			"substitution": mapping.GetSubstitution(),
			"ldap_query":   mapping.GetLdapQuery(),
		}
	}
	return ret
}
