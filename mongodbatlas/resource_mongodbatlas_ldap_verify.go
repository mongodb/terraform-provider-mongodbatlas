package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

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
				ForceNew: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bind_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"authz_query_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"links": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"SUCCESS", "FAILED"},
		Refresh:    resourceLDAPGetStatusRefreshFunc(projectID, ldap.RequestID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
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
	if err != nil || ldapResp == nil {
		return fmt.Errorf(errorLDAPVerifyRead, d.Id(), err)
	}

	if err := d.Set("hostname", ldapResp.Request.Hostname); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "hostname", d.Id(), err)
	}
	if err := d.Set("port", ldapResp.Request.Port); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "port", d.Id(), err)
	}
	if err := d.Set("bind_username", ldapResp.Request.BindUsername); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "bind_username", d.Id(), err)
	}
	if err := d.Set("links", flattenLinks(ldapResp.Links)); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "links", d.Id(), err)
	}
	if err := d.Set("validations", flattenValidations(ldapResp.Validations)); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "validations", d.Id(), err)
	}
	if err := d.Set("request_id", ldapResp.RequestID); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "request_id", d.Id(), err)
	}
	if err := d.Set("status", ldapResp.Status); err != nil {
		return fmt.Errorf(errorLDAPVerifySetting, "status", d.Id(), err)
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

func resourceLDAPGetStatusRefreshFunc(projectID, requestID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.LDAPConfigurations.GetStatus(context.Background(), projectID, requestID)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		return p, p.Status, nil
	}
}
