package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorLDAPVerifyCreate  = "error creating MongoDB LDAPVerify (%s): %s"
	errorLDAPVerifyRead    = "error reading MongoDB LDAPVerify (%s): %s"
	errorLDAPVerifySetting = "error setting `%s` for LDAPVerify(%s): %s"
)

func resourceMongoDBAtlasLDAPVerify() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasLDAPVerifyCreate,
		Read:   resourceMongoDBAtlasLDAPVerifyRead,
		Delete: resourceMongoDBAtlasLDAPVerifyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasLDAPVerifyImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Default:  636,
				Required: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_password": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceMongoDBAtlasLDAPVerifyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	ldapReq := &matlas.LDAP{}

	if v, ok := d.GetOk("hostname"); ok {
		ldapReq.Hostname = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		ldapReq.Port = v.(int)
	}
	if v, ok := d.GetOk("bind_username"); ok {
		ldapReq.BindUsername = v.(string)
	}
	if v, ok := d.GetOk("bind_password"); ok {
		ldapReq.BindPassword = v.(string)
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		ldapReq.CaCertificate = v.(string)
	}
	if v, ok := d.GetOk("authz_query_template"); ok {
		ldapReq.AuthzQueryTemplate = v.(string)
	}

	ldap, _, err := conn.LDAPConfigurations.Verify(context.Background(), projectID, ldapReq)
	if err != nil {
		return fmt.Errorf(errorLDAPVerifyCreate, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldap.RequestID,
	}))

	return resourceMongoDBAtlasLDAPVerifyRead(d, meta)
}

func resourceMongoDBAtlasLDAPVerifyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	requestID := ids["request_id"]

	ldapResp, _, err := conn.LDAPConfigurations.GetStatus(context.Background(), projectID, requestID)
	if err != nil {
		return fmt.Errorf(errorLDAPVerifyRead, d.Id(), err)
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

	return nil
}

func resourceMongoDBAtlasLDAPVerifyDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}

func flattenLinks(linksArray []*matlas.Link) []map[string]interface{} {
	links := make([]map[string]interface{}, 0)
	for _, v := range linksArray {
		links = append(links, map[string]interface{}{
			"href": v.Href,
			"rel":  v.Rel,
		})
	}

	return links
}

func flattenValidations(validationsArray []*matlas.LDAPValidation) []map[string]interface{} {
	validations := make([]map[string]interface{}, 0)
	for _, v := range validationsArray {
		validations = append(validations, map[string]interface{}{
			"status":          v.Status,
			"validation_type": v.ValidationType,
		})
	}

	return validations
}

func resourceMongoDBAtlasLDAPVerifyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a LDAP Verify use the format {project_id}-{request_id}")
	}

	projectID := parts[0]
	requestID := parts[1]

	_, _, err := conn.LDAPConfigurations.GetStatus(context.Background(), projectID, requestID)
	if err != nil {
		return nil, fmt.Errorf(errorLDAPVerifyRead, requestID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorLDAPVerifySetting, "project_id", requestID, err)
	}

	if err := d.Set("request_id", requestID); err != nil {
		return nil, fmt.Errorf(errorLDAPVerifySetting, "request_id", requestID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": requestID,
	}))

	return []*schema.ResourceData{d}, nil
}
