package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasLDAPConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasLDAPConfigurationRead,
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

func dataSourceMongoDBAtlasLDAPConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	ldap, _, err := conn.LDAPConfigurations.Get(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf(errorLDAPConfigurationRead, projectID, err)
	}

	if err = d.Set("authentication_enabled", ldap.LDAP.AuthenticationEnabled); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "authentication_enabled", d.Id(), err)
	}
	if err = d.Set("authorization_enabled", ldap.LDAP.AuthorizationEnabled); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "authorization_enabled", d.Id(), err)
	}
	if err = d.Set("hostname", ldap.LDAP.Hostname); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "hostname", d.Id(), err)
	}
	if err = d.Set("port", ldap.LDAP.Port); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "port", d.Id(), err)
	}
	if err = d.Set("bind_username", ldap.LDAP.BindUsername); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "bind_username", d.Id(), err)
	}
	if err = d.Set("bind_password", ldap.LDAP.BindPassword); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "bind_password", d.Id(), err)
	}
	if err = d.Set("ca_certificate", ldap.LDAP.CaCertificate); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "ca_certificate", d.Id(), err)
	}
	if err = d.Set("authz_query_template", ldap.LDAP.AuthzQueryTemplate); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "authz_query_template", d.Id(), err)
	}
	if err = d.Set("user_to_dn_mapping", flattenDNMapping(ldap.LDAP.UserToDNMapping)); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "user_to_dn_mapping", d.Id(), err)
	}

	d.SetId(projectID)

	return nil
}
