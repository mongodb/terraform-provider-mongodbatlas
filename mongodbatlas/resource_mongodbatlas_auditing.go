package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorAuditingCreate = "error creating MongoDB Auditing (%s): %s"
	errorAuditingUpdate = "error updating MongoDB Auditing (%s): %s"
	errorAuditingRead   = "error reading MongoDB Auditing (%s): %s"
)

func resourceMongoDBAtlasAuditing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasAuditingCreate,
		ReadContext:   resourceMongoDBAtlasAuditingRead,
		UpdateContext: resourceMongoDBAtlasAuditingUpdate,
		DeleteContext: resourceMongoDBAtlasAuditingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audit_authorization_success": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"audit_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"configuration_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasAuditingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	auditingReq := &matlas.Auditing{}

	if auditAuth, ok := d.GetOk("audit_authorization_success"); ok {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool(auditAuth.(bool))
	}

	if auditFilter, ok := d.GetOk("audit_filter"); ok {
		auditingReq.AuditFilter = auditFilter.(string)
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		auditingReq.Enabled = pointy.Bool(enabled.(bool))
	}

	_, _, err := conn.Auditing.Configure(ctx, projectID, auditingReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingCreate, projectID, err))
	}

	d.SetId(projectID)

	return resourceMongoDBAtlasAuditingRead(ctx, d, meta)
}

func resourceMongoDBAtlasAuditingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	auditing, resp, err := conn.Auditing.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("audit_authorization_success", auditing.AuditAuthorizationSuccess); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("audit_filter", auditing.AuditFilter); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("enabled", auditing.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("configuration_type", auditing.ConfigurationType); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasAuditingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection.
	conn := meta.(*MongoDBClient).Atlas

	auditingReq := &matlas.Auditing{}

	if d.HasChange("audit_authorization_success") {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool((d.Get("audit_authorization_success").(bool)))
	}

	if d.HasChange("audit_filter") {
		auditingReq.AuditFilter = d.Get("audit_filter").(string)
	}

	if d.HasChange("enabled") {
		auditingReq.Enabled = pointy.Bool(d.Get("enabled").(bool))
	}

	_, _, err := conn.Auditing.Configure(ctx, d.Id(), auditingReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingUpdate, d.Id(), err))
	}

	return nil
}

func resourceMongoDBAtlasAuditingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
