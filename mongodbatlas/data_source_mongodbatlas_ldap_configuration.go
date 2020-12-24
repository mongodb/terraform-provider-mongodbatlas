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
			"ldap": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

	if err := d.Set("ldap", flattenLDAP(ldap.LDAP)); err != nil {
		return fmt.Errorf(errorLDAPConfigurationSetting, "ldap", projectID, err)
	}

	d.SetId(projectID)

	return nil
}
