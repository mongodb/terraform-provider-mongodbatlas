package auditing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
)

const (
	errorAuditingCreate = "error creating MongoDB Auditing (%s): %s"
	errorAuditingUpdate = "error updating MongoDB Auditing (%s): %s"
	errorAuditingRead   = "error reading MongoDB Auditing (%s): %s"
	errorAuditingDelete = "error deleting MongoDB Auditing (%s): %s"
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
				ForceNew: true,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	auditingReq := &admin.AuditLog{}

	if auditAuth, ok := d.GetOk("audit_authorization_success"); ok {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool(auditAuth.(bool))
	}

	if auditFilter, ok := d.GetOk("audit_filter"); ok {
		auditingReq.AuditFilter = conversion.StringPtr(auditFilter.(string))
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		auditingReq.Enabled = pointy.Bool(enabled.(bool))
	}

	_, _, err := connV2.AuditingApi.UpdateAuditingConfiguration(ctx, projectID, auditingReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingCreate, projectID, err))
	}

	d.SetId(projectID)
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	auditing, resp, err := connV2.AuditingApi.GetAuditingConfiguration(ctx, d.Id()).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("audit_authorization_success", auditing.GetAuditAuthorizationSuccess()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("audit_filter", auditing.GetAuditFilter()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("enabled", auditing.GetEnabled()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	if err := d.Set("configuration_type", auditing.GetConfigurationType()); err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingRead, d.Id(), err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	auditingReq := &admin.AuditLog{}

	if d.HasChange("audit_authorization_success") {
		auditingReq.AuditAuthorizationSuccess = pointy.Bool((d.Get("audit_authorization_success").(bool)))
	}

	if d.HasChange("audit_filter") {
		auditingReq.AuditFilter = conversion.StringPtr(d.Get("audit_filter").(string))
	}

	if d.HasChange("enabled") {
		auditingReq.Enabled = pointy.Bool(d.Get("enabled").(bool))
	}

	_, _, err := connV2.AuditingApi.UpdateAuditingConfiguration(ctx, d.Id(), auditingReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingUpdate, d.Id(), err))
	}

	return nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	auditingReq := &admin.AuditLog{
		Enabled: pointy.Bool(false),
	}
	_, _, err := connV2.AuditingApi.UpdateAuditingConfiguration(ctx, d.Id(), auditingReq).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAuditingDelete, d.Id(), err))
	}
	d.SetId("")
	return nil
}
