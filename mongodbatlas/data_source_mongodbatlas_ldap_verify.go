package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasLDAPVerify() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasLDAPVerifyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Default:  636,
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
			"links": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"validations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasLDAPVerifyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	requestID := d.Get("request_id").(string)

	ldapResp, _, err := conn.LDAPConfigurations.GetStatus(context.Background(), projectID, requestID)
	if err != nil {
		return fmt.Errorf(errorLDAPVerifyRead, projectID, err)
	}

	if err := d.Set("hostname", ldapResp.LDAP.Hostname); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "hostname", d.Id(), err)
	}
	if err := d.Set("port", ldapResp.LDAP.Port); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "port", d.Id(), err)
	}
	if err := d.Set("bind_username", ldapResp.LDAP.BindUsername); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "bind_username", d.Id(), err)
	}
	if err := d.Set("bind_password", ldapResp.LDAP.BindPassword); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "bind_password", d.Id(), err)
	}
	if err := d.Set("ca_certificate", ldapResp.LDAP.CaCertificate); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "ca_certificate", d.Id(), err)
	}
	if err := d.Set("authz_query_template", ldapResp.LDAP.AuthzQueryTemplate); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "authz_query_template", d.Id(), err)
	}
	if err := d.Set("links", flattenLinks(ldapResp.Links)); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "links", d.Id(), err)
	}
	if err := d.Set("validations", flattenValidations(ldapResp.Validations)); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "validations", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldapResp.RequestID,
	}))

	return nil
}
