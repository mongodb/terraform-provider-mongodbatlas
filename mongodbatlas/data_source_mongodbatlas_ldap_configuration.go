package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasLDAPConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasLDAPConfigurationRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"authorization_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bind_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authz_query_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_to_dn_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"substitution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_query": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasLDAPConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	ldap, _, err := conn.LDAPConfigurations.Get(ctx, projectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationRead, projectID, err))
	}

	if err = d.Set("authentication_enabled", ldap.LDAP.AuthenticationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authentication_enabled", d.Id(), err))
	}
	if err = d.Set("authorization_enabled", ldap.LDAP.AuthorizationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authorization_enabled", d.Id(), err))
	}
	if err = d.Set("hostname", ldap.LDAP.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "hostname", d.Id(), err))
	}
	if err = d.Set("port", ldap.LDAP.Port); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "port", d.Id(), err))
	}
	if err = d.Set("bind_username", ldap.LDAP.BindUsername); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "bind_username", d.Id(), err))
	}
	if err = d.Set("bind_password", ldap.LDAP.BindPassword); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "bind_password", d.Id(), err))
	}
	if err = d.Set("ca_certificate", ldap.LDAP.CaCertificate); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "ca_certificate", d.Id(), err))
	}
	if err = d.Set("authz_query_template", ldap.LDAP.AuthzQueryTemplate); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "authz_query_template", d.Id(), err))
	}
	if err = d.Set("user_to_dn_mapping", flattenDNMapping(ldap.LDAP.UserToDNMapping)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPConfigurationSetting, "user_to_dn_mapping", d.Id(), err))
	}

	d.SetId(projectID)

	return nil
}
