package ldapverify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorLDAPVerifyCreate  = "error creating MongoDB LDAPVerify (%s): %s"
	errorLDAPVerifyRead    = "error reading MongoDB LDAPVerify (%s): %s"
	errorLDAPVerifySetting = "error setting `%s` for LDAPVerify(%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasLDAPVerifyCreate,
		ReadContext:   resourceMongoDBAtlasLDAPVerifyRead,
		DeleteContext: resourceMongoDBAtlasLDAPVerifyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasLDAPVerifyImportState,
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

func resourceMongoDBAtlasLDAPVerifyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	ldapReq := &matlas.LDAP{}

	if v, ok := d.GetOk("hostname"); ok {
		ldapReq.Hostname = conversion.Pointer(v.(string))
	}
	if v, ok := d.GetOk("port"); ok {
		ldapReq.Port = conversion.Pointer(v.(int))
	}
	if v, ok := d.GetOk("bind_username"); ok {
		ldapReq.BindUsername = conversion.Pointer(v.(string))
	}
	if v, ok := d.GetOk("bind_password"); ok {
		ldapReq.BindPassword = conversion.Pointer(v.(string))
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		ldapReq.CaCertificate = conversion.Pointer(v.(string))
	}
	if v, ok := d.GetOk("authz_query_template"); ok {
		ldapReq.AuthzQueryTemplate = conversion.Pointer(v.(string))
	}

	ldap, _, err := conn.LDAPConfigurations.Verify(ctx, projectID, ldapReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifyCreate, projectID, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"SUCCESS", "FAILED"},
		Refresh:    resourceLDAPGetStatusRefreshFunc(ctx, projectID, ldap.RequestID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifyCreate, projectID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldap.RequestID,
	}))

	return resourceMongoDBAtlasLDAPVerifyRead(ctx, d, meta)
}

func resourceMongoDBAtlasLDAPVerifyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas

	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	requestID := ids["request_id"]

	ldapResp, resp, err := conn.LDAPConfigurations.GetStatus(context.Background(), projectID, requestID)
	if err != nil || ldapResp == nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorLDAPVerifyRead, d.Id(), err))
	}

	if err := d.Set("hostname", ldapResp.Request.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "hostname", d.Id(), err))
	}
	if err := d.Set("port", ldapResp.Request.Port); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "port", d.Id(), err))
	}
	if err := d.Set("bind_username", ldapResp.Request.BindUsername); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "bind_username", d.Id(), err))
	}
	if err := d.Set("links", FlattenLinks(ldapResp.Links)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "links", d.Id(), err))
	}
	if err := d.Set("validations", flattenValidations(ldapResp.Validations)); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "validations", d.Id(), err))
	}
	if err := d.Set("request_id", ldapResp.RequestID); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "request_id", d.Id(), err))
	}
	if err := d.Set("status", ldapResp.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorLDAPVerifySetting, "status", d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasLDAPVerifyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")

	return nil
}

func FlattenLinks(linksArray []*matlas.Link) []map[string]any {
	links := make([]map[string]any, 0)
	for _, v := range linksArray {
		links = append(links, map[string]any{
			"href": v.Href,
			"rel":  v.Rel,
		})
	}

	return links
}

func flattenValidations(validationsArray []*matlas.LDAPValidation) []map[string]any {
	validations := make([]map[string]any, 0)
	for _, v := range validationsArray {
		validations = append(validations, map[string]any{
			"status":          v.Status,
			"validation_type": v.ValidationType,
		})
	}

	return validations
}

func resourceMongoDBAtlasLDAPVerifyImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a LDAP Verify use the format {project_id}-{request_id}")
	}

	projectID := parts[0]
	requestID := parts[1]

	_, _, err := conn.LDAPConfigurations.GetStatus(ctx, projectID, requestID)
	if err != nil {
		return nil, fmt.Errorf(errorLDAPVerifyRead, requestID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorLDAPVerifySetting, "project_id", requestID, err)
	}

	if err := d.Set("request_id", requestID); err != nil {
		return nil, fmt.Errorf(errorLDAPVerifySetting, "request_id", requestID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": requestID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceLDAPGetStatusRefreshFunc(ctx context.Context, projectID, requestID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		p, resp, err := client.LDAPConfigurations.GetStatus(ctx, projectID, requestID)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		return p, p.Status, nil
	}
}
