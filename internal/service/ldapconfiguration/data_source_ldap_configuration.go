package ldapconfiguration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	resp, _, err := connV2.LDAPConfigurationApi.GetUserSecurity(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRead, projectID, err))
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
	if err = d.Set("bind_password", resp.Ldap.GetBindPassword()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "bind_password", d.Id(), err))
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

	d.SetId(projectID)

	return nil
}
